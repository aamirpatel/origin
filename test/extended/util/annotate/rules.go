package main

// Rules defined here are additive to the rules already defined for
// kube e2e tests in openshift/kubernetes. The kube rules are
// vendored via the following file:
//
//   vendor/k8s.io/kubernetes/openshift-hack/e2e/annotate/rules.go
//
// Rules that are needed to pass the upstream e2e test suite in a
// "default OCP CI" configuration (eg, AWS or GCP, openshift-sdn) must
// be added to openshift/kubernetes to allow CI to pass there, and
// then vendored back into origin. Rules that only apply to
// "non-default" configurations (other clouds, other network
// providers) should be added here.

var (
	testMaps = map[string][]string{
		// tests that require a local host
		"[Local]": {
			// Doesn't work on scaled up clusters
			`\[Feature:ImagePrune\]`,
		},
		// alpha features that are not gated
		"[Disabled:Alpha]": {},
		// tests for features that are not implemented in openshift
		"[Disabled:Unimplemented]": {},
		// tests that rely on special configuration that we do not yet support
		"[Disabled:SpecialConfig]": {},
		// tests that are known broken and need to be fixed upstream or in openshift
		// always add an issue here
		"[Disabled:Broken]": {
			`should idle the service and DeploymentConfig properly`,       // idling with a single service and DeploymentConfig
			`should answer endpoint and wildcard queries for the cluster`, // currently not supported by dns operator https://github.com/openshift/cluster-dns-operator/issues/43

			// https://bugzilla.redhat.com/show_bug.cgi?id=1988272
			`\[sig-network\] Networking should provide Internet connection for containers \[Feature:Networking-IPv6\]`,
			`\[sig-network\] Networking should provider Internet connection for containers using DNS`,

			// https://bugzilla.redhat.com/show_bug.cgi?id=1908645
			`\[sig-network\] Networking Granular Checks: Services should function for service endpoints using hostNetwork`,
			`\[sig-network\] Networking Granular Checks: Services should function for pod-Service\(hostNetwork\)`,

			// https://bugzilla.redhat.com/show_bug.cgi?id=1952460
			`\[sig-network\] Firewall rule control plane should not expose well-known ports`,

			// https://bugzilla.redhat.com/show_bug.cgi?id=1952457
			`\[sig-node\] crictl should be able to run crictl on the node`,

			// https://bugzilla.redhat.com/show_bug.cgi?id=1945091
			`\[Feature:GenericEphemeralVolume\]`,

			// https://bugzilla.redhat.com/show_bug.cgi?id=1953478
			`\[sig-storage\] Dynamic Provisioning Invalid AWS KMS key should report an error and create no PV`,

			// https://bugzilla.redhat.com/show_bug.cgi?id=1957894
			`\[sig-node\] Container Runtime blackbox test when running a container with a new image should be able to pull from private registry with secret`,

			// The new NetworkPolicy test suite is extremely resource
			// intensive and causes itself and other concurrently-running
			// tests to be flaky.
			// https://bugzilla.redhat.com/show_bug.cgi?id=1980141
			`\[sig-network\] Netpol `,

			// https://bugzilla.redhat.com/show_bug.cgi?id=1996128
			`\[sig-network\] \[Feature:IPv6DualStack\] should have ipv4 and ipv6 node podCIDRs`,

			// https://bugzilla.redhat.com/show_bug.cgi?id=2004074
			`\[sig-network-edge\]\[Feature:Idling\] Unidling \[apigroup:apps.openshift.io\]\[apigroup:route.openshift.io\] should work with TCP \(while idling\)`,

			// https://bugzilla.redhat.com/show_bug.cgi?id=2070929
			`\[sig-network\]\[Feature:EgressIP\]\[apigroup:config.openshift.io\] \[internal-targets\]`,

			// https://issues.redhat.com/browse/OCPBUGS-967
			`\[sig-network\] IngressClass \[Feature:Ingress\] should prevent Ingress creation if more than 1 IngressClass marked as default`,

			// https://issues.redhat.com/browse/OCPBUGS-3339
			`\[sig-devex\]\[Feature:ImageEcosystem\]\[mysql\]\[Slow\] openshift mysql image Creating from a template should instantiate the template`,
			`\[sig-devex\]\[Feature:ImageEcosystem\]\[mariadb\]\[Slow\] openshift mariadb image Creating from a template should instantiate the template`,
		},
		// tests that may work, but we don't support them
		"[Disabled:Unsupported]": {
			// Skip vSphere-specific storage tests. The standard in-tree storage tests for vSphere
			// (prefixed with `In-tree Volumes [Driver: vsphere]`) are enough for testing this plugin.
			// https://bugzilla.redhat.com/show_bug.cgi?id=2019115
			`\[sig-storage\].*\[Feature:vsphere\]`,
			// Also, our CI doesn't support topology, so disable those tests
			`\[sig-storage\] In-tree Volumes \[Driver: vsphere\] \[Testpattern: Dynamic PV \(delayed binding\)\] topology should fail to schedule a pod which has topologies that conflict with AllowedTopologies`,
			`\[sig-storage\] In-tree Volumes \[Driver: vsphere\] \[Testpattern: Dynamic PV \(delayed binding\)\] topology should provision a volume and schedule a pod with AllowedTopologies`,
			`\[sig-storage\] In-tree Volumes \[Driver: vsphere\] \[Testpattern: Dynamic PV \(immediate binding\)\] topology should fail to schedule a pod which has topologies that conflict with AllowedTopologies`,
			`\[sig-storage\] In-tree Volumes \[Driver: vsphere\] \[Testpattern: Dynamic PV \(immediate binding\)\] topology should provision a volume and schedule a pod with AllowedTopologies`,
			// Skip openstack-specific storage tests in preparation for in-tree cinder provisioner removal
			// coming with k8s 1.26. This will have to be reverted once 1.26 rebase is effective.
			// https://issues.redhat.com/browse/OCPBUGS-5029
			`\[sig-storage\].*\[Driver: cinder\]`,
		},
		// tests too slow to be part of conformance
		"[Slow]": {},
		// tests that are known flaky
		"[Flaky]": {
			`openshift mongodb replication creating from a template`, // flaking on deployment
		},
		// tests that must be run without competition
		"[Serial]": {
			`\[sig-network\]\[Feature:EgressIP\]`,
		},
		// tests that can't be run in parallel with a copy of itself
		"[Serial:Self]": {
			`\[sig-network\] HostPort validates that there is no conflict between pods with same hostPort but different hostIP and protocol`,
		},
		"[Skipped:azure]": {},
		"[Skipped:ovirt]": {},
		"[Skipped:gce]":   {},

		// These tests are skipped when openshift-tests needs to use a proxy to reach the
		// cluster -- either because the test won't work while proxied, or because the test
		// itself is testing a functionality using it's own proxy.
		"[Skipped:Proxy]": {
			// These tests setup their own proxy, which won't work when we need to access the
			// cluster through a proxy.
			`\[sig-cli\] Kubectl client Simple pod should support exec through an HTTP proxy`,
			`\[sig-cli\] Kubectl client Simple pod should support exec through kubectl proxy`,

			// Kube currently uses the x/net/websockets pkg, which doesn't work with proxies.
			// See: https://github.com/kubernetes/kubernetes/pull/103595
			`\[sig-node\] Pods should support retrieving logs from the container over websockets`,
			`\[sig-cli\] Kubectl Port forwarding With a server listening on localhost should support forwarding over websockets`,
			`\[sig-cli\] Kubectl Port forwarding With a server listening on 0.0.0.0 should support forwarding over websockets`,
			`\[sig-node\] Pods should support remote command execution over websockets`,

			// These tests are flacky and require internet access
			// See https://bugzilla.redhat.com/show_bug.cgi?id=2019375
			`\[sig-builds\]\[Feature:Builds\] build can reference a cluster service with a build being created from new-build should be able to run a build that references a cluster service`,
			`\[sig-builds\]\[Feature:Builds\] oc new-app should succeed with a --name of 58 characters`,
			`\[sig-network\] DNS should resolve DNS of partial qualified names for services`,
			`\[sig-arch\] Only known images used by tests`,
			`\[sig-network\] DNS should provide DNS for the cluster`,
			// This test does not work when using in-proxy cluster, see https://bugzilla.redhat.com/show_bug.cgi?id=2084560
			`\[sig-network\] Networking should provide Internet connection for containers`,
		},
		"[Skipped:SingleReplicaTopology]": {
			`\[sig-apps\] Daemon set \[Serial\] should rollback without unnecessary restarts \[Conformance\]`,
			`\[sig-node\] NoExecuteTaintManager Single Pod \[Serial\] doesn't evict pod with tolerations from tainted nodes`,
			`\[sig-node\] NoExecuteTaintManager Single Pod \[Serial\] eventually evict pod with finite tolerations from tainted nodes`,
			`\[sig-node\] NoExecuteTaintManager Single Pod \[Serial\] evicts pods from tainted nodes`,
			`\[sig-node\] NoExecuteTaintManager Single Pod \[Serial\] removing taint cancels eviction \[Disruptive\] \[Conformance\]`,
			`\[sig-node\] NoExecuteTaintManager Multiple Pods \[Serial\] evicts pods with minTolerationSeconds \[Disruptive\] \[Conformance\]`,
			`\[sig-node\] NoExecuteTaintManager Multiple Pods \[Serial\] only evicts pods without tolerations from tainted nodes`,
			`\[sig-cli\] Kubectl client Kubectl taint \[Serial\] should remove all the taints with the same key off a node`,
		},

		"[Feature:Networking-IPv4]": {
			`\[sig-network\]\[Feature:Router\]\[apigroup:route.openshift.io\] when FIPS is disabled the HAProxy router should serve routes when configured with a 1024-bit RSA key`,
		},

		// Tests that don't pass on disconnected, either due to requiring
		// internet access for GitHub (e.g. many of the s2i builds), or
		// because of pullthrough not supporting ICSP (https://bugzilla.redhat.com/show_bug.cgi?id=1918376)
		"[Skipped:Disconnected]": {
			// Internet access required
			`\[sig-builds\]\[Feature:Builds\] clone repository using git:// protocol should clone using git:// if no proxy is configured`,
			`\[sig-builds\]\[Feature:Builds\] result image should have proper labels set S2I build from a template should create a image from "test-s2i-build.json" template with proper Docker labels`,
			`\[sig-builds\]\[Feature:Builds\] s2i build with a quota Building from a template should create an s2i build with a quota and run it`,
			`\[sig-builds\]\[Feature:Builds\] s2i build with a root user image should create a root build and pass with a privileged SCC`,
			`\[sig-builds\]\[Feature:Builds\]\[timing\] capture build stages and durations should record build stages and durations for docker`,
			`\[sig-builds\]\[Feature:Builds\]\[timing\] capture build stages and durations should record build stages and durations for s2i`,
			`\[sig-builds\]\[Feature:Builds\]\[valueFrom\] process valueFrom in build strategy environment variables should successfully resolve valueFrom in s2i build environment variables`,
			`\[sig-builds\]\[Feature:Builds\]\[volumes\] should mount given secrets and configmaps into the build pod for source strategy builds`,
			`\[sig-builds\]\[Feature:Builds\]\[volumes\] should mount given secrets and configmaps into the build pod for docker strategy builds`,
			`\[sig-builds\]\[Feature:Builds\]\[pullsearch\] docker build where the registry is not specified Building from a Dockerfile whose FROM image ref does not specify the image registry should create a docker build that has buildah search from our predefined list of image registries and succeed`,
			`\[sig-cli\] oc debug ensure it works with image streams`,
			`\[sig-cli\] oc builds complex build start-build`,
			`\[sig-cli\] oc builds complex build webhooks CRUD`,
			`\[sig-cli\] oc builds new-build`,
			`\[sig-devex\] check registry.redhat.io is available and samples operator can import sample imagestreams run sample related validations`,
			`\[sig-devex\]\[Feature:Templates\] templateinstance readiness test should report failed soon after an annotated objects has failed`,
			`\[sig-devex\]\[Feature:Templates\] templateinstance readiness test should report ready soon after all annotated objects are ready`,
			`\[sig-operator\] an end user can use OLM can subscribe to the operator`,
			`\[sig-network\] Networking should provide Internet connection for containers`,
			`\[sig-imageregistry\]\[Serial\] Image signature workflow can push a signed image to openshift registry and verify it`,

			// Need to access non-cached images like ruby and mongodb
			`\[sig-apps\]\[Feature:DeploymentConfig\] deploymentconfigs with multiple image change triggers should run a successful deployment with a trigger used by different containers`,
			`\[sig-apps\]\[Feature:DeploymentConfig\] deploymentconfigs with multiple image change triggers should run a successful deployment with multiple triggers`,

			// ICSP
			`\[sig-apps\]\[Feature:DeploymentConfig\] deploymentconfigs should adhere to Three Laws of Controllers`,
			`\[sig-apps\]\[Feature:DeploymentConfig\] deploymentconfigs adoption will orphan all RCs and adopt them back when recreated`,
			`\[sig-apps\]\[Feature:DeploymentConfig\] deploymentconfigs generation should deploy based on a status version bump`,
			`\[sig-apps\]\[Feature:DeploymentConfig\] deploymentconfigs keep the deployer pod invariant valid should deal with cancellation after deployer pod succeeded`,
			`\[sig-apps\]\[Feature:DeploymentConfig\] deploymentconfigs paused should disable actions on deployments`,
			`\[sig-apps\]\[Feature:DeploymentConfig\] deploymentconfigs rolled back should rollback to an older deployment`,
			`\[sig-apps\]\[Feature:DeploymentConfig\] deploymentconfigs should respect image stream tag reference policy resolve the image pull spec`,
			`\[sig-apps\]\[Feature:DeploymentConfig\] deploymentconfigs viewing rollout history should print the rollout history`,
			`\[sig-apps\]\[Feature:DeploymentConfig\] deploymentconfigs when changing image change trigger should successfully trigger from an updated image`,
			`\[sig-apps\]\[Feature:DeploymentConfig\] deploymentconfigs when run iteratively should only deploy the last deployment`,
			`\[sig-apps\]\[Feature:DeploymentConfig\] deploymentconfigs when tagging images should successfully tag the deployed image`,
			`\[sig-apps\]\[Feature:DeploymentConfig\] deploymentconfigs with custom deployments should run the custom deployment steps`,
			`\[sig-apps\]\[Feature:DeploymentConfig\] deploymentconfigs with enhanced status should include various info in status`,
			`\[sig-apps\]\[Feature:DeploymentConfig\] deploymentconfigs with env in params referencing the configmap should expand the config map key to a value`,
			`\[sig-apps\]\[Feature:DeploymentConfig\] deploymentconfigs with failing hook should get all logs from retried hooks`,
			`\[sig-apps\]\[Feature:DeploymentConfig\] deploymentconfigs with minimum ready seconds set should not transition the deployment to Complete before satisfied`,
			`\[sig-apps\]\[Feature:DeploymentConfig\] deploymentconfigs with revision history limits should never persist more old deployments than acceptable after being observed by the controller`,
			`\[sig-apps\]\[Feature:DeploymentConfig\] deploymentconfigs with test deployments should run a deployment to completion and then scale to zero`,
			`\[sig-apps\]\[Feature:DeploymentConfig\] deploymentconfigs won't deploy RC with unresolved images when patched with empty image`,
			`\[sig-apps\]\[Feature:Jobs\] Users should be able to create and run a job in a user project`,
			`\[sig-arch\] Managed cluster should expose cluster services outside the cluster`,
			`\[sig-arch\]\[Early\] Managed cluster should \[apigroup:config.openshift.io\] start all core operators`,
			`\[sig-auth\]\[Feature:SecurityContextConstraints\] TestPodDefaultCapabilities`,
			`\[sig-builds\]\[Feature:Builds\] Multi-stage image builds should succeed`,
			`\[sig-builds\]\[Feature:Builds\] Optimized image builds should succeed`,
			`\[sig-builds\]\[Feature:Builds\] build can reference a cluster service with a build being created from new-build should be able to run a build that references a cluster service`,
			`\[sig-builds\]\[Feature:Builds\] build have source revision metadata started build should contain source revision information`,
			`\[sig-builds\]\[Feature:Builds\] build with empty source started build should build even with an empty source in build config`,
			`\[sig-builds\]\[Feature:Builds\] build without output image building from templates should create an image from a S2i template without an output image reference defined`,
			`\[sig-builds\]\[Feature:Builds\] build without output image building from templates should create an image from a docker template without an output image reference defined`,
			`\[sig-builds\]\[Feature:Builds\] custom build with buildah being created from new-build should complete build with custom builder image`,
			`\[sig-builds\]\[Feature:Builds\] imagechangetriggers imagechangetriggers should trigger builds of all types`,
			`\[sig-builds\]\[Feature:Builds\] oc new-app should fail with a --name longer than 58 characters`,
			`\[sig-builds\]\[Feature:Builds\] oc new-app should succeed with a --name of 58 characters`,
			`\[sig-builds\]\[Feature:Builds\] oc new-app should succeed with an imagestream`,
			`\[sig-builds\]\[Feature:Builds\] prune builds based on settings in the buildconfig buildconfigs should have a default history limit set when created via the group api`,
			`\[sig-builds\]\[Feature:Builds\] prune builds based on settings in the buildconfig should prune builds after a buildConfig change`,
			`\[sig-builds\]\[Feature:Builds\] prune builds based on settings in the buildconfig should prune canceled builds based on the failedBuildsHistoryLimit setting`,
			`\[sig-builds\]\[Feature:Builds\] prune builds based on settings in the buildconfig should prune completed builds based on the successfulBuildsHistoryLimit setting`,
			`\[sig-builds\]\[Feature:Builds\] prune builds based on settings in the buildconfig should prune errored builds based on the failedBuildsHistoryLimit setting`,
			`\[sig-builds\]\[Feature:Builds\] prune builds based on settings in the buildconfig should prune failed builds based on the failedBuildsHistoryLimit setting`,
			`\[sig-builds\]\[Feature:Builds\] result image should have proper labels set Docker build from a template should create a image from "test-docker-build.json" template with proper Docker labels`,
			`\[sig-builds\]\[Feature:Builds\] verify /run filesystem contents are writeable using a simple Docker Strategy Build`,
			`\[sig-builds\]\[Feature:Builds\] verify /run filesystem contents do not have unexpected content using a simple Docker Strategy Build`,
			`\[sig-builds\]\[Feature:Builds\]\[pullsecret\] docker build using a pull secret Building from a template should create a docker build that pulls using a secret run it`,
			`\[sig-builds\]\[Feature:Builds\]\[valueFrom\] process valueFrom in build strategy environment variables should fail resolving unresolvable valueFrom in docker build environment variable references`,
			`\[sig-builds\]\[Feature:Builds\]\[valueFrom\] process valueFrom in build strategy environment variables should fail resolving unresolvable valueFrom in sti build environment variable references`,
			`\[sig-builds\]\[Feature:Builds\]\[valueFrom\] process valueFrom in build strategy environment variables should successfully resolve valueFrom in docker build environment variables`,
			`\[sig-builds\]\[Feature:Builds\]\[pullsearch\] docker build where the registry is not specified Building from a Dockerfile whose FROM image ref does not specify the image registry should create a docker build that has buildah search from our predefined list of image registries and succeed`,
			`\[sig-cli\] CLI can run inside of a busybox container`,
			`\[sig-cli\] oc debug deployment configs from a build`,
			`\[sig-cli\] oc rsh specific flags should work well when access to a remote shell`,
			`\[sig-cli\] oc builds get buildconfig`,
			`\[sig-cli\] oc builds patch buildconfig`,
			`\[sig-cluster-lifecycle\] Pods cannot access the /config/master API endpoint`,
			`\[sig-imageregistry\]\[Feature:ImageAppend\] Image append should create images by appending them`,
			`\[sig-imageregistry\]\[Feature:ImageExtract\] Image extract should extract content from an image`,
			`\[sig-imageregistry\]\[Feature:ImageInfo\] Image info should display information about images`,
			`\[sig-imageregistry\]\[Feature:ImageLayers\] Image layer subresource should return layers from tagged images`,
			`\[sig-imageregistry\]\[Feature:ImageTriggers\] Annotation trigger reconciles after the image is overwritten`,
			`\[sig-imageregistry\]\[Feature:Image\] oc tag should change image reference for internal images`,
			`\[sig-imageregistry\]\[Feature:Image\] oc tag should work when only imagestreams api is available`,
			`\[sig-instrumentation\] Prometheus \[apigroup:image.openshift.io\] when installed on the cluster should have a AlertmanagerReceiversNotConfigured alert in firing state`,
			`\[sig-instrumentation\] Prometheus \[apigroup:image.openshift.io\] when installed on the cluster should have important platform topology metrics`,
			`\[sig-instrumentation\] Prometheus \[apigroup:image.openshift.io\] when installed on the cluster should have non-Pod host cAdvisor metrics`,
			`\[sig-instrumentation\] Prometheus \[apigroup:image.openshift.io\] when installed on the cluster should provide ingress metrics`,
			`\[sig-instrumentation\] Prometheus \[apigroup:image.openshift.io\] when installed on the cluster should provide named network metrics`,
			`\[sig-instrumentation\] Prometheus \[apigroup:image.openshift.io\] when installed on the cluster should report telemetry \[Late\]`,
			`\[sig-instrumentation\] Prometheus \[apigroup:image.openshift.io\] when installed on the cluster should start and expose a secured proxy and unsecured metrics`,
			`\[sig-instrumentation\] Prometheus \[apigroup:image.openshift.io\] when installed on the cluster shouldn't have failing rules evaluation`,
			`\[sig-instrumentation\] Prometheus \[apigroup:image.openshift.io\] when installed on the cluster shouldn't report any alerts in firing state apart from Watchdog and AlertmanagerReceiversNotConfigured \[Early\]`,
			`\[sig-instrumentation\] Prometheus \[apigroup:image.openshift.io\] when installed on the cluster when using openshift-sdn should be able to get the sdn ovs flows`,
			`\[sig-instrumentation\]\[Late\] OpenShift alerting rules \[apigroup:image.openshift.io\] should have a valid severity label`,
			`\[sig-instrumentation\]\[Late\] OpenShift alerting rules \[apigroup:image.openshift.io\] should have description and summary annotations`,
			`\[sig-instrumentation\]\[Late\] OpenShift alerting rules \[apigroup:image.openshift.io\] should have a runbook_url annotation if the alert is critical`,
			`\[sig-instrumentation\]\[Late\] Alerts should have a Watchdog alert in firing state the entire cluster run`,
			`\[sig-instrumentation\]\[Late\] Alerts shouldn't exceed the 500 series limit of total series sent via telemetry from each cluster`,
			`\[sig-instrumentation\]\[Late\] Alerts shouldn't report any alerts in firing or pending state apart from Watchdog and AlertmanagerReceiversNotConfigured and have no gaps in Watchdog firing`,
			`\[sig-instrumentation\]\[sig-builds\]\[Feature:Builds\] Prometheus when installed on the cluster should start and expose a secured proxy and verify build metrics`,
			`\[sig-network-edge\]\[Conformance\]\[Area:Networking\]\[Feature:Router\] The HAProxy router should be able to connect to a service that is idled because a GET on the route will unidle it`,
			`\[sig-network\]\[Feature:Router\]\[apigroup:config.openshift.io\] The HAProxy router should enable openshift-monitoring to pull metrics`,
			`\[sig-network\]\[Feature:Router\]\[apigroup:config.openshift.io\] The HAProxy router should expose a health check on the metrics port`,
			`\[sig-network\]\[Feature:Router\]\[apigroup:config.openshift.io\] The HAProxy router should expose prometheus metrics for a route`,
			`\[sig-network\]\[Feature:Router\]\[apigroup:config.openshift.io\] The HAProxy router should expose the profiling endpoints`,
			`\[sig-network\]\[Feature:Router\]\[apigroup:route.openshift.io\] The HAProxy router should override the route host for overridden domains with a custom value`,
			`\[sig-network\]\[Feature:Router\]\[apigroup:route.openshift.io\] The HAProxy router should override the route host with a custom value`,
			`\[sig-network\]\[Feature:Router\]\[apigroup:operator.openshift.io\] The HAProxy router should respond with 503 to unrecognized hosts`,
			`\[sig-network\]\[Feature:Router\]\[apigroup:route.openshift.io\] The HAProxy router should run even if it has no access to update status`,
			`\[sig-network\]\[Feature:Router\]\[apigroup:image.openshift.io\] The HAProxy router should serve a route that points to two services and respect weights`,
			`\[sig-network\]\[Feature:Router\]\[apigroup:operator.openshift.io\] The HAProxy router should serve routes that were created from an ingress`,
			`\[sig-network\]\[Feature:Router\]\[apigroup:route.openshift.io\] The HAProxy router should serve the correct routes when scoped to a single namespace and label set`,
			`\[sig-network\]\[Feature:Router\]\[apigroup:config.openshift.io\]\[apigroup:operator.openshift.io\] The HAProxy router should set Forwarded headers appropriately`,
			`\[sig-network\]\[Feature:Router\]\[apigroup:route.openshift.io\]\[apigroup:operator.openshift.io\] The HAProxy router should support reencrypt to services backed by a serving certificate automatically`,
			`\[sig-network\] Networking should provide Internet connection for containers \[Feature:Networking-IPv6\]`,
			`\[sig-node\] Managed cluster should report ready nodes the entire duration of the test run`,
			`\[sig-storage\]\[Late\] Metrics should report short attach times`,
			`\[sig-storage\]\[Late\] Metrics should report short mount times`,
		},

		// tests that don't pass under openshift-sdn NetworkPolicy mode are specified
		// in the rules file in openshift/kubernetes, not here.

		// tests that don't pass under openshift-sdn multitenant mode
		"[Skipped:Network/OpenShiftSDN/Multitenant]": {
			`\[Feature:NetworkPolicy\]`, // not compatible with multitenant mode
		},
		// tests that don't pass under OVN Kubernetes
		"[Skipped:Network/OVNKubernetes]": {
			// ovn-kubernetes does not support named ports
			`NetworkPolicy.*named port`,
		},
		"[Skipped:ibmroks]": {
			// skip Gluster tests (not supported on ROKS worker nodes)
			// https://bugzilla.redhat.com/show_bug.cgi?id=1825009 - e2e: skip Glusterfs-related tests upstream for rhel7 worker nodes
			`\[Driver: gluster\]`,
			`GlusterFS`,
			`GlusterDynamicProvisioner`,

			// Nodes in ROKS have access to secrets in the cluster to handle encryption
			// https://bugzilla.redhat.com/show_bug.cgi?id=1825013 - ROKS: worker nodes have access to secrets in the cluster
			`\[sig-auth\] \[Feature:NodeAuthorizer\] Getting a non-existent configmap should exit with the Forbidden error, not a NotFound error`,
			`\[sig-auth\] \[Feature:NodeAuthorizer\] Getting a non-existent secret should exit with the Forbidden error, not a NotFound error`,
			`\[sig-auth\] \[Feature:NodeAuthorizer\] Getting a secret for a workload the node has access to should succeed`,
			`\[sig-auth\] \[Feature:NodeAuthorizer\] Getting an existing configmap should exit with the Forbidden error`,
			`\[sig-auth\] \[Feature:NodeAuthorizer\] Getting an existing secret should exit with the Forbidden error`,

			// Access to node external address is blocked from pods within a ROKS cluster by Calico
			// https://bugzilla.redhat.com/show_bug.cgi?id=1825016 - e2e: NodeAuthenticator tests use both external and internal addresses for node
			`\[sig-auth\] \[Feature:NodeAuthenticator\] The kubelet's main port 10250 should reject requests with no credentials`,
			`\[sig-auth\] \[Feature:NodeAuthenticator\] The kubelet can delegate ServiceAccount tokens to the API server`,

			// Calico is allowing the request to timeout instead of returning 'REFUSED'
			// https://bugzilla.redhat.com/show_bug.cgi?id=1825021 - ROKS: calico SDN results in a request timeout when accessing services with no endpoints
			`\[sig-network\] Services should be rejected when no endpoints exist`,

			// Mode returned by RHEL7 worker contains an extra character not expected by the test: dgtrwx vs dtrwx
			// https://bugzilla.redhat.com/show_bug.cgi?id=1825024 - e2e: Failing test - HostPath should give a volume the correct mode
			`\[sig-storage\] HostPath should give a volume the correct mode`,

			// Currently ibm-master-proxy-static and imbcloud-block-storage-plugin tolerate all taints
			// https://bugzilla.redhat.com/show_bug.cgi?id=1825027
			`\[Feature:Platform\] Managed cluster should ensure control plane operators do not make themselves unevictable`,
		},
		// Tests which can't be run/don't make sense to run against a cluster with all optional capabilities disabled
		"[Skipped:NoOptionalCapabilities]": {
			// Most storage tests don't pass when the storage capability is disabled.
			// this list needs to be refined as there are some storage tests we should be able to run.
			// Tracker for enabling more storage tests: https://issues.redhat.com/browse/OCPPLAN-9509
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Dynamic PV \(block volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv used in a pod that is deleted while the kubelet is down cleans up when the kubelet returns.`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Dynamic PV \(block volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv used in a pod that is force deleted while the kubelet is down cleans up when the kubelet returns.`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Dynamic PV \(block volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv written before kubelet restart is readable after restart.`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Dynamic PV \(block volmode\)\] multiVolume \[Slow\] should concurrently access the volume and restored snapshot from pods on the same node`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Dynamic PV \(block volmode\)\] provisioning should provision storage with snapshot data source`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Dynamic PV \(default fs\)\] provisioning should provision storage with snapshot data source`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should unmount if pod is force deleted while kubelet is down`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should unmount if pod is gracefully deleted while kubelet is down`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Dynamic PV \(filesystem volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv written before kubelet restart is readable after restart.`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Dynamic PV \(filesystem volmode\)\] volume-lifecycle-performance should provision volumes at scale within performance constraints`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Dynamic Snapshot \(delete policy\)\] snapshottable-stress\[Feature:VolumeSnapshotDataSource\] should support snapshotting of many volumes repeatedly`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Dynamic Snapshot \(delete policy\)\] snapshottable\[Feature:VolumeSnapshotDataSource\] volume snapshot controller should check snapshot fields, check restore correctly works after modifying source data, check deletion (persistent)`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Dynamic Snapshot \(retain policy\)\] snapshottable-stress\[Feature:VolumeSnapshotDataSource\] should support snapshotting of many volumes repeatedly`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Dynamic Snapshot \(retain policy\)\] snapshottable\[Feature:VolumeSnapshotDataSource\] volume snapshot controller should check snapshot fields, check restore correctly works after modifying source data, check deletion (persistent)`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Ephemeral Snapshot \(delete policy\)\] snapshottable\[Feature:VolumeSnapshotDataSource\] volume snapshot controller should check snapshot fields, check restore correctly works, check deletion (ephemeral)`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Ephemeral Snapshot \(retain policy\)\] snapshottable\[Feature:VolumeSnapshotDataSource\] volume snapshot controller should check snapshot fields, check restore correctly works, check deletion (ephemeral)`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Pre-provisioned Snapshot \(delete policy\)\] snapshottable\[Feature:VolumeSnapshotDataSource\] volume snapshot controller should check snapshot fields, check restore correctly works after modifying source data, check deletion (persistent)`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Pre-provisioned Snapshot \(retain policy\)\] snapshottable\[Feature:VolumeSnapshotDataSource\] volume snapshot controller should check snapshot fields, check restore correctly works after modifying source data, check deletion (persistent)`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Dynamic Snapshot \(delete policy\)\] snapshottable-stress[Feature:VolumeSnapshotDataSource] should support snapshotting of many volumes repeatedly`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Dynamic Snapshot \(delete policy\)\] snapshottable[Feature:VolumeSnapshotDataSource] volume snapshot controller should check snapshot fields, check restore correctly works after modifying source data, check deletion (persistent)`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Dynamic Snapshot \(retain policy\)\] snapshottable-stress[Feature:VolumeSnapshotDataSource] should support snapshotting of many volumes repeatedly`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Dynamic Snapshot \(retain policy\)\] snapshottable[Feature:VolumeSnapshotDataSource] volume snapshot controller should check snapshot fields, check restore correctly works after modifying source data, check deletion (persistent)`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Ephemeral Snapshot \(delete policy\)\] snapshottable[Feature:VolumeSnapshotDataSource] volume snapshot controller should check snapshot fields, check restore correctly works, check deletion (ephemeral)`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Ephemeral Snapshot \(retain policy\)\] snapshottable[Feature:VolumeSnapshotDataSource] volume snapshot controller should check snapshot fields, check restore correctly works, check deletion (ephemeral)`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Pre-provisioned Snapshot \(delete policy\)\] snapshottable[Feature:VolumeSnapshotDataSource] volume snapshot controller should check snapshot fields, check restore correctly works after modifying source data, check deletion (persistent)`,
			`\[sig-storage\] CSI Volumes \[Driver: csi-hostpath\] \[Testpattern: Pre-provisioned Snapshot \(retain policy\)\] snapshottable[Feature:VolumeSnapshotDataSource] volume snapshot controller should check snapshot fields, check restore correctly works after modifying source data, check deletion (persistent)`,

			`\[sig-storage\] CSI mock volume CSI Volume Snapshots \[Feature:VolumeSnapshotDataSource\] volumesnapshotcontent and pvc in Bound state with deletion timestamp set should not get deleted while snapshot finalizer exists`,
			`\[sig-storage\] CSI mock volume CSI Volume Snapshots secrets \[Feature:VolumeSnapshotDataSource\] volume snapshot create/delete with secrets`,

			`\[sig-storage\] Dynamic Provisioning DynamicProvisioner \[Slow\] \[Feature:StorageProvider\] deletion should be idempotent`,
			`\[sig-storage\] Dynamic Provisioning DynamicProvisioner \[Slow\] \[Feature:StorageProvider\] should provision storage with different parameters`,

			`\[sig-storage\] HostPathType Block Device \[Slow\] Should be able to mount block device \'ablkdev\' successfully when HostPathType is HostPathBlockDev`,
			`\[sig-storage\] HostPathType Block Device \[Slow\] Should be able to mount block device \'ablkdev\' successfully when HostPathType is HostPathUnset`,
			`\[sig-storage\] HostPathType Block Device \[Slow\] Should fail on mounting block device \'ablkdev\' when HostPathType is HostPathCharDev`,
			`\[sig-storage\] HostPathType Block Device \[Slow\] Should fail on mounting block device \'ablkdev\' when HostPathType is HostPathDirectory`,
			`\[sig-storage\] HostPathType Block Device \[Slow\] Should fail on mounting block device \'ablkdev\' when HostPathType is HostPathFile`,
			`\[sig-storage\] HostPathType Block Device \[Slow\] Should fail on mounting block device \'ablkdev\' when HostPathType is HostPathSocket`,
			`\[sig-storage\] HostPathType Block Device \[Slow\] Should fail on mounting non-existent block device \'does-not-exist-blk-dev\' when HostPathType is HostPathBlockDev`,
			`\[sig-storage\] HostPathType Character Device \[Slow\] Should be able to mount character device \'achardev\' successfully when HostPathType is HostPathCharDev`,
			`\[sig-storage\] HostPathType Character Device \[Slow\] Should be able to mount character device \'achardev\' successfully when HostPathType is HostPathUnset`,
			`\[sig-storage\] HostPathType Character Device \[Slow\] Should fail on mounting character device \'achardev\' when HostPathType is HostPathBlockDev`,
			`\[sig-storage\] HostPathType Character Device \[Slow\] Should fail on mounting character device \'achardev\' when HostPathType is HostPathDirectory`,
			`\[sig-storage\] HostPathType Character Device \[Slow\] Should fail on mounting character device \'achardev\' when HostPathType is HostPathFile`,
			`\[sig-storage\] HostPathType Character Device \[Slow\] Should fail on mounting character device \'achardev\' when HostPathType is HostPathSocket`,
			`\[sig-storage\] HostPathType Character Device \[Slow\] Should fail on mounting non-existent character device \'does-not-exist-char-dev\' when HostPathType is HostPathCharDev`,
			`\[sig-storage\] HostPathType Socket \[Slow\] Should be able to mount socket \'asocket\' successfully when HostPathType is HostPathSocket`,
			`\[sig-storage\] HostPathType Socket \[Slow\] Should be able to mount socket \'asocket\' successfully when HostPathType is HostPathUnset`,
			`\[sig-storage\] HostPathType Socket \[Slow\] Should fail on mounting non-existent socket \'does-not-exist-socket\' when HostPathType is HostPathSocket`,
			`\[sig-storage\] HostPathType Socket \[Slow\] Should fail on mounting socket \'asocket\' when HostPathType is HostPathBlockDev`,
			`\[sig-storage\] HostPathType Socket \[Slow\] Should fail on mounting socket \'asocket\' when HostPathType is HostPathCharDev`,
			`\[sig-storage\] HostPathType Socket \[Slow\] Should fail on mounting socket \'asocket\' when HostPathType is HostPathDirectory`,
			`\[sig-storage\] HostPathType Socket \[Slow\] Should fail on mounting socket \'asocket\' when HostPathType is HostPathFile`,

			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(block volmode\)\(allowExpansion\)\] volume-expand should resize volume when PVC is edited while pod is using it`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Generic Ephemeral-volume \(default fs\) \(immediate-binding\)\] ephemeral should support two pods which have the same volume definition`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should fail if subpath file is outside the volume`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should support creating multiple subpath from same volumes`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Inline-volume \(xfs\)\]\[Slow\] volumes should store data`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(delayed binding\)\] topology should provision a volume and schedule a pod with AllowedTopologies`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(block volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv written before kubelet restart is readable after restart.`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should support readOnly directory specified in the volumeMount`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should support readOnly file specified in the volumeMount`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Pre-provisioned PV \(filesystem volmode\)\] volumeMode should not mount / map unused volumes in a pod`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(block volmode\)\] volumeMode should not mount / map unused volumes in a pod`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should support existing directories when readOnly specified in the volumeSource`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(xfs\)\]\[Slow\] multiVolume \[Slow\] should access to two volumes with different volume mode and retain data across pod recreation on the same node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] provisioning should provision storage with mount options`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(xfs\)\]\[Slow\] multiVolume \[Slow\] should concurrently access the single read-only volume from pods on the same node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should verify container cannot write to subpath readonly volumes`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(filesystem volmode\)\] multiVolume \[Slow\] should concurrently access the single volume from pods on the same node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(ext4\)\] multiVolume \[Slow\] should concurrently access the single read-only volume from pods on the same node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(block volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv used in a pod that is deleted while the kubelet is down cleans up when the kubelet returns.`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(ext4\)\] multiVolume \[Slow\] should concurrently access the single volume from pods on the same node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Pre-provisioned PV \(ext4\)\] volumes should store data`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should support existing single file`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] fsgroupchangepolicy \(Always\)\[LinuxOnly\], pod created with an initial fsgroup, volume contents ownership changed via chgrp in first pod, new pod with different fsgroup applied to the volume contents`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] fsgroupchangepolicy \(Always\)\[LinuxOnly\], pod created with an initial fsgroup, volume contents ownership changed via chgrp in first pod, new pod with same fsgroup applied to the volume contents`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] fsgroupchangepolicy \(Always\)\[LinuxOnly\], pod created with an initial fsgroup, new pod fsgroup applied to volume contents`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should support restarting containers using file as subpath`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(immediate binding\)\] topology should provision a volume and schedule a pod with AllowedTopologies`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(block volmode\)\] multiVolume \[Slow\] should access to two volumes with the same volume mode and retain data across pod recreation on different node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Generic Ephemeral-volume \(default fs\) \(late-binding\)\] ephemeral should support multiple inline ephemeral volumes`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(block volmode\)\] multiVolume \[Slow\] should concurrently access the single volume from pods on the same node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] volumes should allow exec of files on the volume`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Generic Ephemeral-volume \(default fs\) \(late-binding\)\] ephemeral should support two pods which have the same volume definition`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should support restarting containers using directory as subpath`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(ext4\)\] multiVolume \[Slow\] should access to two volumes with different volume mode and retain data across pod recreation on different node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should fail if subpath directory is outside the volume`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Generic Ephemeral-volume \(default fs\) \(immediate-binding\)\] ephemeral should support expansion of pvcs created for ephemeral pvcs`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(block volmode\)\] multiVolume \[Slow\] should access to two volumes with different volume mode and retain data across pod recreation on the same node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Pre-provisioned PV \(filesystem volmode\)\] volumeMode should fail to use a volume in a pod with mismatched mode`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(ext4\)\] volumes should allow exec of files on the volume`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Generic Ephemeral-volume \(default fs\) \(late-binding\)\] ephemeral should create read/write inline ephemeral volume`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Generic Ephemeral-volume \(block volmode\) \(late-binding\)\] ephemeral should support multiple inline ephemeral volumes`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Pre-provisioned PV \(xfs\)\]\[Slow\] volumes should store data`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(ext4\)\] multiVolume \[Slow\] should access to two volumes with different volume mode and retain data across pod recreation on the same node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Inline-volume \(ext4\)\] volumes should store data`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should be able to unmount after the subpath directory is deleted`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] volumes should store data`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(filesystem volmode\)\] multiVolume \[Slow\] should access to two volumes with the same volume mode and retain data across pod recreation on different node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(block volmode\)\] multiVolume \[Slow\] should access to two volumes with the same volume mode and retain data across pod recreation on the same node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] provisioning should mount multiple PV pointing to the same storage on the same node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Pre-provisioned PV \(block volmode\)\] volumeMode should fail to use a volume in a pod with mismatched mode`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Pre-provisioned PV \(default fs\)\] volumes should allow exec of files on the volume`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(xfs\)\]\[Slow\] multiVolume \[Slow\] should access to two volumes with different volume mode and retain data across pod recreation on different node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(block volmode\)\] multiVolume \[Slow\] should access to two volumes with different volume mode and retain data across pod recreation on different node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(ext4\)\] volumes should store data`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(ext4\)\] multiVolume \[Slow\] should access to two volumes with the same volume mode and retain data across pod recreation on the same node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Inline-volume \(default fs\)\] volumes should allow exec of files on the volume`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(xfs\)\]\[Slow\] volumes should store data`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should support non-existent path`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Pre-provisioned PV \(block volmode\)\] volumeMode should not mount / map unused volumes in a pod`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Inline-volume \(default fs\)\] volumes should store data`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should unmount if pod is gracefully deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(xfs\)\]\[Slow\] multiVolume \[Slow\] should access to two volumes with the same volume mode and retain data across pod recreation on the same node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(xfs\)\]\[Slow\] volumes should allow exec of files on the volume`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(filesystem volmode\)\] multiVolume \[Slow\] should concurrently access the single read-only volume from pods on the same node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should fail if subpath with backstepping is outside the volume`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Inline-volume \(ext4\)\] volumes should allow exec of files on the volume`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(block volmode\)\(allowExpansion\)\] volume-expand Verify if offline PVC expansion works`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(block volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv used in a pod that is force deleted while the kubelet is down cleans up when the kubelet returns.`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(block volmode\)\] multiVolume \[Slow\] should concurrently access the single read-only volume from pods on the same node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(filesystem volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv written before kubelet restart is readable after restart.`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Generic Ephemeral-volume \(default fs\) \(late-binding\)\] ephemeral should create read-only inline ephemeral volume`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Pre-provisioned PV \(default fs\)\] volumes should store data`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should support existing directory`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Inline-volume \(xfs\)\]\[Slow\] volumes should allow exec of files on the volume`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] volumeIO should write files of various sizes, verify size, validate content`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Generic Ephemeral-volume \(default fs\) \(immediate-binding\)\] ephemeral should create read-only inline ephemeral volume`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should unmount if pod is force deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Generic Ephemeral-volume \(default fs\) \(late-binding\)\] ephemeral should support expansion of pvcs created for ephemeral pvcs`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] fsgroupchangepolicy \(OnRootMismatch\)\[LinuxOnly\], pod created with an initial fsgroup, volume contents ownership changed via chgrp in first pod, new pod with same fsgroup skips ownership changes to the volume contents`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Pre-provisioned PV \(xfs\)\]\[Slow\] volumes should allow exec of files on the volume`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\(allowExpansion\)\] volume-expand Verify if offline PVC expansion works`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] fsgroupchangepolicy \(OnRootMismatch\)\[LinuxOnly\], pod created with an initial fsgroup, new pod fsgroup applied to volume contents`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Pre-provisioned PV \(block volmode\)\] volumes should store data`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\(allowExpansion\)\] volume-expand should resize volume when PVC is edited while pod is using it`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should support file as subpath`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(block volmode\)\] volumes should store data`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Generic Ephemeral-volume \(block volmode\) \(late-binding\)\] ephemeral should create read/write inline ephemeral volume`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] fsgroupchangepolicy \(OnRootMismatch\)\[LinuxOnly\], pod created with an initial fsgroup, volume contents ownership changed via chgrp in first pod, new pod with different fsgroup applied to the volume contents`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Pre-provisioned PV \(ext4\)\] volumes should allow exec of files on the volume`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(xfs\)\]\[Slow\] multiVolume \[Slow\] should access to two volumes with the same volume mode and retain data across pod recreation on different node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(ext4\)\] multiVolume \[Slow\] should access to two volumes with the same volume mode and retain data across pod recreation on different node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Generic Ephemeral-volume \(default fs\) \(immediate-binding\)\] ephemeral should create read/write inline ephemeral volume`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should fail if non-existent subpath is outside the volume`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(filesystem volmode\)\] volumeMode should not mount / map unused volumes in a pod`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Generic Ephemeral-volume \(block volmode\) \(late-binding\)\] ephemeral should support two pods which have the same volume definition`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(filesystem volmode\)\] multiVolume \[Slow\] should access to two volumes with the same volume mode and retain data across pod recreation on the same node`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Generic Ephemeral-volume \(block volmode\) \(late-binding\)\] ephemeral should support expansion of pvcs created for ephemeral pvcs`,
			`\[sig-storage\] In-tree Volumes \[Driver: aws\] \[Testpattern: Dynamic PV \(xfs\)\]\[Slow\] multiVolume \[Slow\] should concurrently access the single volume from pods on the same node`,

			`\[sig-storage\] In-tree Volumes \[Driver: ceph\]\[Feature:Volumes\]\[Serial\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should fail if subpath file is outside the volume`,
			`\[sig-storage\] In-tree Volumes \[Driver: ceph\]\[Feature:Volumes\]\[Serial\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should unmount if pod is force deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: ceph\]\[Feature:Volumes\]\[Serial\] \[Testpattern: Inline-volume \(default fs\)\] subPath should unmount if pod is force deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: ceph\]\[Feature:Volumes\]\[Serial\] \[Testpattern: Pre-provisioned PV \(default fs\)\] volumes should store data`,
			`\[sig-storage\] In-tree Volumes \[Driver: ceph\]\[Feature:Volumes\]\[Serial\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should verify container cannot write to subpath readonly volumes`,
			`\[sig-storage\] In-tree Volumes \[Driver: ceph\]\[Feature:Volumes\]\[Serial\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should unmount if pod is gracefully deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: ceph\]\[Feature:Volumes\]\[Serial\] \[Testpattern: Inline-volume \(default fs\)\] subPath should unmount if pod is gracefully deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: ceph\]\[Feature:Volumes\]\[Serial\] \[Testpattern: Pre-provisioned PV \(filesystem volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv written before kubelet restart is readable after restart.`,
			`\[sig-storage\] In-tree Volumes \[Driver: ceph\]\[Feature:Volumes\]\[Serial\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should support readOnly directory specified in the volumeMount`,

			`\[sig-storage\] In-tree Volumes \[Driver: emptydir\] \[Testpattern: Inline-volume \(default fs\)\] subPath should unmount if pod is gracefully deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: emptydir\] \[Testpattern: Inline-volume \(default fs\)\] subPath should unmount if pod is force deleted while kubelet is down`,

			`\[sig-storage\] In-tree Volumes \[Driver: gluster\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should unmount if pod is force deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: gluster\] \[Testpattern: Pre-provisioned PV \(filesystem volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv written before kubelet restart is readable after restart.`,
			`\[sig-storage\] In-tree Volumes \[Driver: gluster\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should unmount if pod is gracefully deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: gluster\] \[Testpattern: Inline-volume \(default fs\)\] subPath should unmount if pod is gracefully deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: gluster\] \[Testpattern: Inline-volume \(default fs\)\] subPath should unmount if pod is force deleted while kubelet is down`,

			`\[sig-storage\] In-tree Volumes \[Driver: hostPathSymlink\] \[Testpattern: Inline-volume \(default fs\)\] volumeIO should write files of various sizes, verify size, validate content`,

			`\[sig-storage\] In-tree Volumes \[Driver: hostPath\] \[Testpattern: Inline-volume \(default fs\)\] volumeIO should write files of various sizes, verify size, validate content`,

			`\[sig-storage\] In-tree Volumes \[Driver: iscsi\]\[Feature:Volumes\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should unmount if pod is force deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: iscsi\]\[Feature:Volumes\] \[Testpattern: Pre-provisioned PV \(filesystem volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv written before kubelet restart is readable after restart.`,
			`\[sig-storage\] In-tree Volumes \[Driver: iscsi\]\[Feature:Volumes\] \[Testpattern: Pre-provisioned PV \(block volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv used in a pod that is force deleted while the kubelet is down cleans up when the kubelet returns.`,
			`\[sig-storage\] In-tree Volumes \[Driver: iscsi\]\[Feature:Volumes\] \[Testpattern: Pre-provisioned PV \(block volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv used in a pod that is deleted while the kubelet is down cleans up when the kubelet returns.`,
			`\[sig-storage\] In-tree Volumes \[Driver: iscsi\]\[Feature:Volumes\] \[Testpattern: Pre-provisioned PV \(block volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv written before kubelet restart is readable after restart.`,
			`\[sig-storage\] In-tree Volumes \[Driver: iscsi\]\[Feature:Volumes\] \[Testpattern: Inline-volume \(default fs\)\] subPath should unmount if pod is force deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: iscsi\]\[Feature:Volumes\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should unmount if pod is gracefully deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: iscsi\]\[Feature:Volumes\] \[Testpattern: Inline-volume \(default fs\)\] subPath should unmount if pod is gracefully deleted while kubelet is down`,

			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: block\] \[Testpattern: Pre-provisioned PV \(block volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv used in a pod that is deleted while the kubelet is down cleans up when the kubelet returns.`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: block\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should unmount if pod is force deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: block\] \[Testpattern: Pre-provisioned PV \(block volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv used in a pod that is force deleted while the kubelet is down cleans up when the kubelet returns.`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: block\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should unmount if pod is gracefully deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: dir-link\] \[Testpattern: Pre-provisioned PV \(filesystem volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv written before kubelet restart is readable after restart.`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: block\] \[Testpattern: Pre-provisioned PV \(filesystem volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv written before kubelet restart is readable after restart.`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: blockfs\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should unmount if pod is gracefully deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: blockfs\] \[Testpattern: Pre-provisioned PV \(filesystem volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv written before kubelet restart is readable after restart.`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: tmpfs\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should unmount if pod is force deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: dir-link\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should unmount if pod is gracefully deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: dir-link-bindmounted\] \[Testpattern: Pre-provisioned PV \(filesystem volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv written before kubelet restart is readable after restart.`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: dir\] \[Testpattern: Pre-provisioned PV \(filesystem volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv written before kubelet restart is readable after restart.`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: dir-bindmounted\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should unmount if pod is gracefully deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: dir\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should unmount if pod is force deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: dir-bindmounted\] \[Testpattern: Pre-provisioned PV \(filesystem volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv written before kubelet restart is readable after restart.`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: dir-link-bindmounted\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should unmount if pod is gracefully deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: dir-bindmounted\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should unmount if pod is force deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: tmpfs\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should unmount if pod is gracefully deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: dir\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should unmount if pod is gracefully deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: dir-link\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should unmount if pod is force deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: blockfs\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should unmount if pod is force deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: dir-link-bindmounted\] \[Testpattern: Pre-provisioned PV \(default fs\)\] subPath should unmount if pod is force deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: block\] \[Testpattern: Pre-provisioned PV \(block volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv written before kubelet restart is readable after restart.`,
			`\[sig-storage\] In-tree Volumes \[Driver: local\]\[LocalVolumeType: tmpfs\] \[Testpattern: Pre-provisioned PV \(filesystem volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv written before kubelet restart is readable after restart.`,

			`\[sig-storage\] In-tree Volumes \[Driver: nfs\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should unmount if pod is gracefully deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: nfs\] \[Testpattern: Dynamic PV \(default fs\)\] subPath should unmount if pod is force deleted while kubelet is down`,
			`\[sig-storage\] In-tree Volumes \[Driver: nfs\] \[Testpattern: Dynamic PV \(filesystem volmode\)\] disruptive\[Disruptive\]\[LinuxOnly\] Should test that pv written before kubelet restart is readable after restart.`,

			// This test requires a valid console url which doesn't exist when the optional console capability is disabled.
			`\[sig-cli\] oc basics can show correct whoami result with console`,
		},
	}
)
