// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package api

import "github.com/Azure/aks-engine-azurestack/pkg/api/common"

// staticSchedulerConfig is not user-overridable
var staticSchedulerConfig = map[string]string{
	"--kubeconfig":   "/var/lib/kubelet/kubeconfig",
	"--leader-elect": "true",
}

// defaultSchedulerConfig provides targeted defaults, but is user-overridable
var defaultSchedulerConfig = map[string]string{
	"--v":               "2",
	"--profiling":       DefaultKubernetesSchedulerEnableProfiling,
	"--bind-address":    "127.0.0.1",    // STIG Rule ID: SV-242384r879530_rule
	"--tls-min-version": "VersionTLS12", // STIG Rule ID: SV-242377r879519_rule
}

func (cs *ContainerService) setSchedulerConfig() {
	o := cs.Properties.OrchestratorProfile

	// If no user-configurable scheduler config values exists, make an empty map, and fill in with defaults
	if o.KubernetesConfig.SchedulerConfig == nil {
		o.KubernetesConfig.SchedulerConfig = make(map[string]string)
	}

	for key, val := range defaultSchedulerConfig {
		// If we don't have a user-configurable scheduler config for each option
		if _, ok := o.KubernetesConfig.SchedulerConfig[key]; !ok {
			// then assign the default value
			o.KubernetesConfig.SchedulerConfig[key] = val
		}
	}

	// STIG Rule ID: SV-254801r879719_rule
	addDefaultFeatureGates(o.KubernetesConfig.SchedulerConfig, o.OrchestratorVersion, "1.25.0", "PodSecurity=true")

	// We don't support user-configurable values for the following,
	// so any of the value assignments below will override user-provided values
	for key, val := range staticSchedulerConfig {
		o.KubernetesConfig.SchedulerConfig[key] = val
	}

	invalidFeatureGates := []string{}
	// Remove --feature-gate VolumeSnapshotDataSource starting with 1.22
	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.22.0-alpha.1") {
		invalidFeatureGates = append(invalidFeatureGates, "VolumeSnapshotDataSource")
	}
	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.27.0") {
		// Remove --feature-gate ControllerManagerLeaderMigration starting with 1.27
		// Reference: https://github.com/kubernetes/kubernetes/pull/113534
		invalidFeatureGates = append(invalidFeatureGates, "ControllerManagerLeaderMigration")
		// Remove --feature-gate ExpandCSIVolumes, ExpandInUsePersistentVolumes, ExpandPersistentVolumes starting with 1.27
		// Reference: https://github.com/kubernetes/kubernetes/pull/113942
		invalidFeatureGates = append(invalidFeatureGates, "ExpandCSIVolumes", "ExpandInUsePersistentVolumes", "ExpandPersistentVolumes")
		// Remove --feature-gate CSIInlineVolume, CSIMigration, CSIMigrationAzureDisk, DaemonSetUpdateSurge, EphemeralContainers, IdentifyPodOS, LocalStorageCapacityIsolation, NetworkPolicyEndPort, StatefulSetMinReadySeconds starting with 1.27
		// Reference: https://github.com/kubernetes/kubernetes/pull/114410
		invalidFeatureGates = append(invalidFeatureGates, "CSIInlineVolume", "CSIMigration", "CSIMigrationAzureDisk", "DaemonSetUpdateSurge", "EphemeralContainers", "IdentifyPodOS", "LocalStorageCapacityIsolation", "NetworkPolicyEndPort", "StatefulSetMinReadySeconds")
	}

	if common.IsKubernetesVersionGe(o.OrchestratorVersion, "1.28.0") {
		// Remove --feature-gate AdvancedAuditing starting with 1.28
		invalidFeatureGates = append(invalidFeatureGates, "AdvancedAuditing","DisableAcceleratorUsageMetrics","DryRun","PodSecurity")

		invalidFeatureGates = append(invalidFeatureGates, "NetworkPolicyStatus","PodHasNetworkCondition","UserNamespacesStatelessPodsSupport")

		// Remove --feature-gate CSIMigrationGCE starting with 1.28
		// Reference: https://github.com/kubernetes/kubernetes/pull/117055
		invalidFeatureGates = append(invalidFeatureGates, "CSIMigrationGCE")

		// Remove --feature-gate CSIStorageCapacity starting with 1.28
		// Reference: https://github.com/kubernetes/kubernetes/pull/118018
		invalidFeatureGates = append(invalidFeatureGates, "CSIStorageCapacity")

		// Remove --feature-gate DelegateFSGroupToCSIDriver starting with 1.28
		// Reference: https://github.com/kubernetes/kubernetes/pull/117655
		invalidFeatureGates = append(invalidFeatureGates, "DelegateFSGroupToCSIDriver")

		// Remove --feature-gate DevicePlugins starting with 1.28
		// Reference: https://github.com/kubernetes/kubernetes/pull/117656
		invalidFeatureGates = append(invalidFeatureGates, "DevicePlugins")

		// Remove --feature-gate KubeletCredentialProviders starting with 1.28
		// Reference: https://github.com/kubernetes/kubernetes/pull/116901
		invalidFeatureGates = append(invalidFeatureGates, "KubeletCredentialProviders")

		// Remove --feature-gate MixedProtocolLBService, ServiceInternalTrafficPolicy, ServiceIPStaticSubrange, EndpointSliceTerminatingCondition  starting with 1.28
		// Reference: https://github.com/kubernetes/kubernetes/pull/117237
		invalidFeatureGates = append(invalidFeatureGates, "MixedProtocolLBService","ServiceInternalTrafficPolicy","ServiceIPStaticSubrange","EndpointSliceTerminatingCondition")

		// Remove --feature-gate WindowsHostProcessContainers starting with 1.28
		// Reference: https://github.com/kubernetes/kubernetes/pull/117570
		invalidFeatureGates = append(invalidFeatureGates, "WindowsHostProcessContainers")
	}
	removeInvalidFeatureGates(o.KubernetesConfig.SchedulerConfig, invalidFeatureGates)
}
