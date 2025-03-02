package monitorapi

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"

	corev1 "k8s.io/api/core/v1"
)

func LocatePod(pod *corev1.Pod) string {
	return fmt.Sprintf("ns/%s pod/%s node/%s uid/%s", pod.Namespace, pod.Name, pod.Spec.NodeName, pod.UID)
}

func LocatePodContainer(pod *corev1.Pod, containerName string) string {
	return fmt.Sprintf("ns/%s pod/%s node/%s uid/%s container/%s", pod.Namespace, pod.Name, pod.Spec.NodeName, pod.UID, containerName)
}

// NonUniquePodLocator produces an inexact locator based on namespace and name.  This is useful when dealing with events
// that are produced that do not contain UIDs.  Ultimately, we should use UIDs everywhere, but this is will keep some our
// matching working until then.
func NonUniquePodLocatorFrom(locator string) string {
	parts := LocatorParts(locator)
	namespace := NamespaceFrom(parts)
	return fmt.Sprintf("ns/%s pod/%s", namespace, parts["pod"])
}

func PodFrom(locator string) PodReference {
	parts := LocatorParts(locator)
	namespace := NamespaceFrom(parts)
	name := parts["pod"]
	uid := parts["uid"]
	if len(namespace) == 0 || len(name) == 0 || len(uid) == 0 {
		return PodReference{}
	}
	return PodReference{
		NamespacedReference: NamespacedReference{
			Namespace: namespace,
			Name:      name,
			UID:       uid,
		},
	}
}

func ContainerFrom(locator string) ContainerReference {
	pod := PodFrom(locator)
	parts := LocatorParts(locator)
	name := parts["container"]
	if len(name) == 0 || len(pod.UID) == 0 {
		return ContainerReference{}
	}
	return ContainerReference{
		Pod:           pod,
		ContainerName: name,
	}
}

type PodReference struct {
	NamespacedReference
}

func (r PodReference) ToLocator() string {
	return fmt.Sprintf("ns/%s pod/%s uid/%s", r.Namespace, r.Name, r.UID)
}

type ContainerReference struct {
	Pod           PodReference
	ContainerName string
}

func (r ContainerReference) ToLocator() string {
	return fmt.Sprintf("ns/%s pod/%s uid/%s container/%s", r.Pod.Namespace, r.Pod.Name, r.Pod.UID, r.ContainerName)
}

func AnnotationsFromMessage(message string) map[string]string {
	tokens := strings.Split(message, " ")
	annotations := map[string]string{}
	for _, curr := range tokens {
		if !strings.Contains(curr, "/") {
			continue
		}
		annotationTokens := strings.Split(curr, "/")
		annotations[annotationTokens[0]] = annotationTokens[1]
	}
	return annotations
}

func ReasonFrom(message string) string {
	annotations := AnnotationsFromMessage(message)
	return annotations["reason"]
}

func PhaseFrom(message string) string {
	annotations := AnnotationsFromMessage(message)
	return annotations["phase"]
}

func ReasonedMessage(reason string, message ...string) string {
	return fmt.Sprintf("reason/%v %s", reason, strings.Join(message, "; "))
}

func ReasonedMessagef(reason, messageFormat string, a ...interface{}) string {
	return ReasonedMessage(reason, fmt.Sprintf(messageFormat, a...))
}

const (
	// PodIPReused means the same pod IP is in use by two pods at the same time.
	PodIPReused = "ReusedPodIP"

	PodReasonCreated               = "Created"
	PodReasonGracefulDeleteStarted = "GracefulDelete"
	PodReasonDeleted               = "Deleted"
	PodReasonScheduled             = "Scheduled"

	ContainerReasonContainerExit      = "ContainerExit"
	ContainerReasonContainerStart     = "ContainerStart"
	ContainerReasonContainerWait      = "ContainerWait"
	ContainerReasonReadinessFailed    = "ReadinessFailed"
	ContainerReasonReadinessErrored   = "ReadinessErrored"
	ContainerReasonStartupProbeFailed = "StartupProbeFailed"
	ContainerReasonReady              = "Ready"
	ContainerReasonNotReady           = "NotReady"

	PodReasonDeletedBeforeScheduling = "DeletedBeforeScheduling"
	PodReasonDeletedAfterCompletion  = "DeletedAfterCompletion"
)

var (
	// PodLifecycleTransitionReasons are the reasons associated with non-overlapping pod lifecycle states.
	// A pod is logically identified by UID (I bet it's a name right now).
	// Pods don't exist before create and don't exist after delete.
	// Between those two states, each of these reasons can be ordered by time and used to create a contiguous view
	// into the lifecycle of a pod.
	PodLifecycleTransitionReasons = sets.NewString(
		PodReasonCreated,
		PodReasonScheduled,
		PodReasonGracefulDeleteStarted,
		PodReasonDeleted,
	)

	// ContainerLifecycleTransitionReasons are the reasons associated with non-overlapping container lifecycle states.
	// The logical beginning and end are based on ContainerWait and ContainerExit.
	// A container is logically identified by a Pod plus a container name.
	ContainerLifecycleTransitionReasons = sets.NewString(
		ContainerReasonContainerWait,
		ContainerReasonContainerStart,
		ContainerReasonContainerExit,
	)

	// ContainerReadinessTransitionReasons are the reasons associated with non-overlapping container readiness states.
	// A container is logically identified by a Pod plus a container name.
	// The logical beginning and end are based on ContainerStart and ContainerExit, with initial state of ready=false and final state of ready=false.
	// Each of these reasons can be ordered by time and used to create a contiguous view into the lifecycle of a pod.
	ContainerReadinessTransitionReasons = sets.NewString(
		ContainerReasonReady,
		ContainerReasonNotReady,
	)

	KubeletReadinessCheckReasons = sets.NewString(
		ContainerReasonReadinessFailed,
		ContainerReasonReadinessErrored,
		ContainerReasonStartupProbeFailed,
	)
)

type ByTimeWithNamespacedPods []EventInterval

func (intervals ByTimeWithNamespacedPods) Less(i, j int) bool {
	lhsIsPodConstructed := strings.Contains(intervals[i].Message, "constructed") && strings.Contains(intervals[i].Locator, "pod/")
	rhsIsPodConstructed := strings.Contains(intervals[j].Message, "constructed") && strings.Contains(intervals[j].Locator, "pod/")
	switch {
	case lhsIsPodConstructed && rhsIsPodConstructed:
		lhsNamespace := NamespaceFromLocator(intervals[i].Locator)
		rhsNamespace := NamespaceFromLocator(intervals[j].Locator)
		if lhsNamespace < rhsNamespace {
			return true
		} else if lhsNamespace > rhsNamespace {
			return false
		} else {
			// sort on time, so fall through.
		}
	case lhsIsPodConstructed && !rhsIsPodConstructed:
		return true
	case !lhsIsPodConstructed && rhsIsPodConstructed:
		return false
	case !lhsIsPodConstructed && !rhsIsPodConstructed:
		// fall through
	}

	switch d := intervals[i].From.Sub(intervals[j].From); {
	case d < 0:
		return true
	case d > 0:
		return false
	}
	switch d := intervals[i].To.Sub(intervals[j].To); {
	case d < 0:
		return true
	case d > 0:
		return false
	}
	return intervals[i].Message < intervals[j].Message
}

func (intervals ByTimeWithNamespacedPods) Len() int { return len(intervals) }
func (intervals ByTimeWithNamespacedPods) Swap(i, j int) {
	intervals[i], intervals[j] = intervals[j], intervals[i]
}
