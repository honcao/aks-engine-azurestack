// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package api

import (
	"testing"

	"github.com/Azure/go-autorest/autorest/to"
)

func TestControllerManagerConfigEnableRbac(t *testing.T) {
	// Test EnableRbac = true
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.EnableRbac = to.BoolPtr(true)
	cs.setControllerManagerConfig()
	cm := cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--use-service-account-credentials"] != "true" {
		t.Fatalf("got unexpected '--use-service-account-credentials' Controller Manager config value for EnableRbac=true: %s",
			cm["--use-service-account-credentials"])
	}

	// Test default
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.EnableRbac = to.BoolPtr(false)
	cs.setControllerManagerConfig()
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--use-service-account-credentials"] != DefaultKubernetesCtrlMgrUseSvcAccountCreds {
		t.Fatalf("got unexpected '--use-service-account-credentials' Controller Manager config value for EnableRbac=false: %s",
			cm["--use-service-account-credentials"])
	}
}

func TestControllerManagerConfigCloudProvider(t *testing.T) {
	// Test UseCloudControllerManager = true
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.UseCloudControllerManager = to.BoolPtr(true)
	cs.setControllerManagerConfig()
	cm := cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--cloud-provider"] != "external" {
		t.Fatalf("got unexpected '--cloud-provider' Controller Manager config value for UseCloudControllerManager=true: %s",
			cm["--cloud-provider"])
	}

	// Test UseCloudControllerManager = false
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.UseCloudControllerManager = to.BoolPtr(false)
	cs.setControllerManagerConfig()
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--cloud-provider"] != "azure" {
		t.Fatalf("got unexpected '--cloud-provider' Controller Manager config value for UseCloudControllerManager=false: %s",
			cm["--cloud-provider"])
	}
}

func TestControllerManagerConfigEnableProfiling(t *testing.T) {
	// Test
	// "controllerManagerConfig": {
	// 	"--profiling": "true"
	// },
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig = map[string]string{
		"--profiling": "true",
	}
	cs.setControllerManagerConfig()
	cm := cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--profiling"] != "true" {
		t.Fatalf("got unexpected '--profiling' Controller Manager config value for \"--profiling\": \"true\": %s",
			cm["--profiling"])
	}

	// Test default
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.setControllerManagerConfig()
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--profiling"] != DefaultKubernetesCtrMgrEnableProfiling {
		t.Fatalf("got unexpected default value for '--profiling' Controller Manager config: %s",
			cm["--profiling"])
	}
}

func TestControllerManagerConfigFeatureGates(t *testing.T) {
	// test defaultTestClusterVer
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.setControllerManagerConfig()
	cm := cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--feature-gates"] != "LegacyServiceAccountTokenNoAutoGeneration=false,LocalStorageCapacityIsolation=true,PodSecurity=true" {
		t.Fatalf("got unexpected '--feature-gates' Controller Manager config value for \"--feature-gates\": \"LegacyServiceAccountTokenNoAutoGeneration=false,LocalStorageCapacityIsolation=true,PodSecurity=true\": %s",
			cm["--feature-gates"])
	}

	// test 1.19.0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.19.0"
	cs.setControllerManagerConfig()
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--feature-gates"] != "LocalStorageCapacityIsolation=true" {
		t.Fatalf("got unexpected '--feature-gates' Controller Manager config value for \"--feature-gates\": \"LocalStorageCapacityIsolation=true\": %s",
			cm["--feature-gates"])
	}

	// test 1.22.0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.22.0"
	cs.setControllerManagerConfig()
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--feature-gates"] != "LocalStorageCapacityIsolation=true" {
		t.Fatalf("got unexpected '--feature-gates' Controller Manager config value for \"--feature-gates\": \"LocalStorageCapacityIsolation=true\": %s",
			cm["--feature-gates"])
	}

	// test 1.24.0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.24.0"
	cs.setControllerManagerConfig()
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--feature-gates"] != "LegacyServiceAccountTokenNoAutoGeneration=false,LocalStorageCapacityIsolation=true" {
		t.Fatalf("got unexpected '--feature-gates' Controller Manager config value for \"--feature-gates\": \"LegacyServiceAccountTokenNoAutoGeneration=false,LocalStorageCapacityIsolation=true\": %s",
			cm["--feature-gates"])
	}

	// test 1.25.0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.25.0"
	cs.setControllerManagerConfig()
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--feature-gates"] != "LegacyServiceAccountTokenNoAutoGeneration=false,LocalStorageCapacityIsolation=true,PodSecurity=true" {
		t.Fatalf("got unexpected '--feature-gates' Controller Manager config value for \"--feature-gates\": \"LegacyServiceAccountTokenNoAutoGeneration=false,LocalStorageCapacityIsolation=true,PodSecurity=true\": %s",
			cm["--feature-gates"])
	}

	// test 1.26.0
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.26.0"
	cs.setControllerManagerConfig()
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--feature-gates"] != "LegacyServiceAccountTokenNoAutoGeneration=false,LocalStorageCapacityIsolation=true,PodSecurity=true" {
		t.Fatalf("got unexpected '--feature-gates' Controller Manager config value for \"--feature-gates\": \"LegacyServiceAccountTokenNoAutoGeneration=false,LocalStorageCapacityIsolation=true,PodSecurity=true\": %s",
			cm["--feature-gates"])
	}

	// test user-overrides
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	cm["--feature-gates"] = "TaintBasedEvictions=true"
	cs.setControllerManagerConfig()
	if cm["--feature-gates"] != "LegacyServiceAccountTokenNoAutoGeneration=false,LocalStorageCapacityIsolation=true,PodSecurity=true,TaintBasedEvictions=true" {
		t.Fatalf("got unexpected '--feature-gates' Controller Manager config value for \"--feature-gates\": \"LocalStorageCapacityIsolation=true,PodSecurity=true,TaintBasedEvictions=true\": %s",
			cm["--feature-gates"])
	}

	// test user-overrides, removal of VolumeSnapshotDataSource for k8s versions >= 1.22
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.22.0"
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	cm["--feature-gates"] = "VolumeSnapshotDataSource=true"
	cs.setControllerManagerConfig()
	if cm["--feature-gates"] != "LocalStorageCapacityIsolation=true" {
		t.Fatalf("got unexpected '--feature-gates' Controller Manager config value for \"--feature-gates\": \"LocalStorageCapacityIsolation=true\": %s",
			cm["--feature-gates"])
	}

	// test user-overrides, no removal of VolumeSnapshotDataSource for k8s versions < 1.22
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.OrchestratorProfile.OrchestratorVersion = "1.19.0"
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	cm["--feature-gates"] = "VolumeSnapshotDataSource=true"
	cs.setControllerManagerConfig()
	if cm["--feature-gates"] != "LocalStorageCapacityIsolation=true,VolumeSnapshotDataSource=true" {
		t.Fatalf("got unexpected '--feature-gates' Controller Manager config value for \"--feature-gates\": \"LocalStorageCapacityIsolation=true,VolumeSnapshotDataSource=true\": %s",
			cm["--feature-gates"])
	}
}

func TestControllerManagerDefaultConfig(t *testing.T) {
	// Azure defaults
	cs := CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.setControllerManagerConfig()
	cm := cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--node-monitor-grace-period"] != string(DefaultKubernetesCtrlMgrNodeMonitorGracePeriod) {
		t.Fatalf("expected controller-manager to have node-monitor-grace-period set to its default value")
	}
	if cm["--pod-eviction-timeout"] != string(DefaultKubernetesCtrlMgrPodEvictionTimeout) {
		t.Fatalf("expected controller-manager to have pod-eviction-timeout set to its default value")
	}
	if cm["--route-reconciliation-period"] != string(DefaultKubernetesCtrlMgrRouteReconciliationPeriod) {
		t.Fatalf("expected controller-manager to have route-reconciliation-period set to its default value")
	}
	if cm["--bind-address"] != "127.0.0.1" {
		t.Fatalf("expected controller-manager to have route-reconciliation-period set to its default value")
	}
	if cm["--tls-min-version"] != "VersionTLS12" {
		t.Fatalf("expected controller-manager to have route-reconciliation-period set to its default value")
	}

	// Azure Stack defaults
	cs = CreateMockContainerService("testcluster", defaultTestClusterVer, 3, 2, false)
	cs.Properties.CustomCloudProfile = &CustomCloudProfile{}
	cs.setControllerManagerConfig()
	cm = cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig
	if cm["--node-monitor-grace-period"] != string(DefaultAzureStackKubernetesCtrlMgrNodeMonitorGracePeriod) {
		t.Fatalf("expected controller-manager to have node-monitor-grace-period set to its default value")
	}
	if cm["--pod-eviction-timeout"] != string(DefaultAzureStackKubernetesCtrlMgrPodEvictionTimeout) {
		t.Fatalf("expected controller-manager to have pod-eviction-timeout set to its default value")
	}
	if cm["--route-reconciliation-period"] != string(DefaultAzureStackKubernetesCtrlMgrRouteReconciliationPeriod) {
		t.Fatalf("expected controller-manager to have route-reconciliation-period set to its default value")
	}
}

func TestControllerManagerInsecureFlag(t *testing.T) {
	type controllerManagerTest struct {
		version string
		found   bool
	}

	controllerManagerTestsForceDelete := []controllerManagerTest{
		{
			version: "1.23.0",
			found:   true,
		},
		{
			version: "1.24.0",
			found:   false,
		},
	}

	for _, tt := range controllerManagerTestsForceDelete {
		cs := CreateMockContainerService("testcluster", tt.version, 3, 2, false)
		cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig = map[string]string{
			"--address": "0.0.0.0",
			"--port":    "443",
		}
		cs.setControllerManagerConfig()
		a := cs.Properties.OrchestratorProfile.KubernetesConfig.ControllerManagerConfig

		_, found := a["--address"]
		if found != tt.found {
			t.Fatalf("got --address found %t want %t", found, tt.found)
		}
		_, found = a["--port"]
		if found != tt.found {
			t.Fatalf("got --port found %t want %t", found, tt.found)
		}
	}

}
