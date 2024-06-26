//
//
// Copyright Red Hat
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"fmt"
	"reflect"
	"testing"

	registryv1alpha1 "github.com/devfile/registry-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIsTLSEnabled(t *testing.T) {
	tlsEnabled := true
	tlsDisabled := false

	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want bool
	}{
		{
			name: "Case 1: TLS enabled in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					TLS: registryv1alpha1.DevfileRegistrySpecTLS{
						Enabled: &tlsEnabled,
					},
				},
			},
			want: true,
		},
		{
			name: "Case 2: TLS disabled in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					TLS: registryv1alpha1.DevfileRegistrySpecTLS{
						Enabled: &tlsDisabled,
					},
				},
			},
			want: false,
		},
		{
			name: "Case 3: TLS not set, default set to true",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tlsSetting := IsTLSEnabled(&tt.cr)
			if tlsSetting != tt.want {
				t.Errorf("TestIsTLSEnabled error: tls value mismatch, expected: %v got: %v", tt.want, tlsSetting)
			}
		})
	}

}

func TestIsStorageEnabled(t *testing.T) {
	storageEnabled := true
	storageDisabled := false

	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want bool
	}{
		{
			name: "Case 1: Storage enabled in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					Storage: registryv1alpha1.DevfileRegistrySpecStorage{
						Enabled: &storageEnabled,
					},
				},
			},
			want: true,
		},
		{
			name: "Case 2: Storage disabled in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					Storage: registryv1alpha1.DevfileRegistrySpecStorage{
						Enabled: &storageDisabled,
					},
				},
			},
			want: false,
		},
		{
			name: "Case 3: Storage not set, default set to false",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tlsSetting := IsStorageEnabled(&tt.cr)
			if tlsSetting != tt.want {
				t.Errorf("TestIsStorageEnabled error: storage value mismatch, expected: %v got: %v", tt.want, tlsSetting)
			}
		})
	}

}

func TestGetDevfileRegistryVolumeSource(t *testing.T) {
	storageEnabled := true
	storageDisabled := false
	crName := "devfileregistry-test"
	crTestData := []registryv1alpha1.DevfileRegistry{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: crName,
			},
			Spec: registryv1alpha1.DevfileRegistrySpec{
				Storage: registryv1alpha1.DevfileRegistrySpecStorage{
					Enabled: &storageEnabled,
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: crName,
			},
			Spec: registryv1alpha1.DevfileRegistrySpec{
				Storage: registryv1alpha1.DevfileRegistrySpecStorage{
					Enabled: &storageDisabled,
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: crName,
			},
			Spec: registryv1alpha1.DevfileRegistrySpec{},
		},
	}

	tests := []struct {
		name string
		cr   *registryv1alpha1.DevfileRegistry
		want corev1.VolumeSource
	}{
		{
			name: "Case 1: Storage enabled in DevfileRegistry CR",
			cr:   &crTestData[0],
			want: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: PVCName(&crTestData[0]),
				},
			},
		},
		{
			name: "Case 2: Storage disabled in DevfileRegistry CR",
			cr:   &crTestData[1],
			want: corev1.VolumeSource{},
		},
		{
			name: "Case 3: Storage not set, default set to false",
			cr:   &crTestData[2],
			want: corev1.VolumeSource{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tlsSetting := GetDevfileRegistryVolumeSource(tt.cr)
			if !reflect.DeepEqual(tlsSetting, tt.want) {
				t.Errorf("TestGetDevfileRegistryVolumeSource error: storage source mismatch, expected: %v got: %v", tt.want, tlsSetting)
			}
		})
	}

}

func TestGetDevfileIndexMemoryLimit(t *testing.T) {
	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want resource.Quantity
	}{
		{
			name: "Case 1: Memory Limit size set in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					DevfileIndex: registryv1alpha1.DevfileRegistrySpecContainer{
						MemoryLimit: "5Gi",
					},
				},
			},
			want: resource.MustParse("5Gi"),
		},
		{
			name: "Case 2:  Memory Limit size not set in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					DevfileIndex: registryv1alpha1.DevfileRegistrySpecContainer{},
				},
			},
			want: resource.MustParse(DefaultDevfileIndexMemoryLimit),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			volSize := GetDevfileIndexMemoryLimit(&tt.cr)
			if volSize != tt.want {
				t.Errorf("TestGetDevfileIndexMemoryLimit error: storage size mismatch, expected: %v got: %v", tt.want, volSize)
			}
		})
	}

}

func TestGetOCIRegistryMemoryLimit(t *testing.T) {
	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want resource.Quantity
	}{
		{
			name: "Case 1: Memory Limit size set in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					OciRegistry: registryv1alpha1.DevfileRegistrySpecContainer{
						MemoryLimit: "5Gi",
					},
				},
			},
			want: resource.MustParse("5Gi"),
		},
		{
			name: "Case 2:  Memory Limit size not set in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					OciRegistry: registryv1alpha1.DevfileRegistrySpecContainer{},
				},
			},
			want: resource.MustParse(DefaultOCIRegistryMemoryLimit),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			volSize := GetOCIRegistryMemoryLimit(&tt.cr)
			if volSize != tt.want {
				t.Errorf("TestGetOCIRegistryMemoryLimit error: storage size mismatch, expected: %v got: %v", tt.want, volSize)
			}
		})
	}

}

func TestGetRegistryViewerMemoryLimit(t *testing.T) {
	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want resource.Quantity
	}{
		{
			name: "Case 1: Memory Limit size set in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					RegistryViewer: registryv1alpha1.DevfileRegistrySpecContainer{
						MemoryLimit: "5Gi",
					},
				},
			},
			want: resource.MustParse("5Gi"),
		},
		{
			name: "Case 2:  Memory Limit size not set in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					RegistryViewer: registryv1alpha1.DevfileRegistrySpecContainer{},
				},
			},
			want: resource.MustParse(DefaultRegistryViewerMemoryLimit),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			volSize := GetRegistryViewerMemoryLimit(&tt.cr)
			if volSize != tt.want {
				t.Errorf("TestGetRegistryViewerMemoryLimit error: storage size mismatch, expected: %v got: %v", tt.want, volSize)
			}
		})
	}

}

func TestGetDevfileRegistryVolumeSize(t *testing.T) {
	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want string
	}{
		{
			name: "Case 1: Volume size set in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					Storage: registryv1alpha1.DevfileRegistrySpecStorage{
						RegistryVolumeSize: "5Gi",
					},
				},
			},
			want: "5Gi",
		},
		{
			name: "Case 2: Volume size not set in DevfileRegistry CR",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{},
			},
			want: DefaultDevfileRegistryVolumeSize,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			volSize := getDevfileRegistryVolumeSize(&tt.cr)
			if volSize != tt.want {
				t.Errorf("TestGetDevfileRegistryVolumeSize error: storage size mismatch, expected: %v got: %v", tt.want, volSize)
			}
		})
	}

}

func TestIsTelemetryEnabled(t *testing.T) {
	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want bool
	}{
		{
			name: "Case 1: Telemetry key not specified",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					Telemetry: registryv1alpha1.DevfileRegistrySpecTelemetry{
						RegistryName: "test",
					},
				},
			},
			want: false,
		},
		{
			name: "Case 2: Telemetry key specified",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					Telemetry: registryv1alpha1.DevfileRegistrySpecTelemetry{
						RegistryName: "test",
						Key:          "abcdef",
					},
				},
			},
			want: true,
		},
		{
			name: "Case 3: Telemetry key specified but is empty",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					Telemetry: registryv1alpha1.DevfileRegistrySpecTelemetry{
						RegistryName: "test",
						Key:          "",
					},
				},
			},
			want: false,
		},
		{
			name: "Case 4: Telemetry object is empty",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enabled := IsTelemetryEnabled(&tt.cr)
			if enabled != tt.want {
				t.Errorf("func TestIsTelemetryEnabled(t *testing.T) {\n error: enablement value mismatch, expected: %v got: %v", tt.want, enabled)
			}
		})
	}

}

func Test_getDevfileRegistrySpecContainer(t *testing.T) {
	tests := []struct {
		name         string
		quantity     string
		defaultValue string
		want         resource.Quantity
	}{
		{
			name:         "Case 1: DevfileRegistrySpecContainer given correct quantity",
			quantity:     "256Mi",
			defaultValue: "512Mi",
			want:         resource.MustParse("256Mi"),
		},
		{
			name:         "Case 2: DevfileRegistrySpecContainer given correct quantity",
			quantity:     "test",
			defaultValue: "512Mi",
			want:         resource.MustParse("512Mi"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getDevfileRegistrySpecContainer(tt.quantity, tt.defaultValue)
			if result != tt.want {
				t.Errorf("func TestgetDevfileRegistrySpecContainer(t *testing.T) {\n error: enablement value mismatch, expected: %v got: %v", tt.want, result)
			}
		})
	}
}

func TestGetK8sIngressClass(t *testing.T) {
	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want string
	}{
		{
			name: "Case 1: K8s ingress class set",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					K8s: registryv1alpha1.DevfileRegistrySpecK8sOnly{
						IngressClass: "test",
					},
				},
			},
			want: "test",
		},
		{
			name: "Case 2: K8s ingress class not set",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					Telemetry: registryv1alpha1.DevfileRegistrySpecTelemetry{},
				},
			},
			want: DefaultK8sIngressClass,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetK8sIngressClass(&tt.cr)
			if result != tt.want {
				t.Errorf("func TestGetK8sIngressClass(t *testing.T) {\n error: enablement value mismatch, expected: %v got: %v", tt.want, result)
			}
		})
	}
}

func TestGetHostnameOverride(t *testing.T) {
	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want string
	}{
		{
			name: "Case 1: Hostname override is set",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					HostnameOverride: "192.168.1.123.nip.io",
				},
			},
			want: "192.168.1.123.nip.io",
		},
		{
			name: "Case 2: Hostname override is unset",
			cr:   registryv1alpha1.DevfileRegistry{},
			want: DefaultHostnameOverride,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetHostnameOverride(&tt.cr)
			if result != tt.want {
				t.Errorf("func TestGetHostnameOverride(t *testing.T) {\n error: enablement value mismatch, expected: %v got: %v", tt.want, result)
			}
		})
	}
}

func TestGetNameOverride(t *testing.T) {
	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want string
	}{
		{
			name: "Case 1: App name override is set",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "devfile-registry-test",
				},
			},
			want: "devfile-registry-test",
		},
		{
			name: "Case 2: App name override is unset",
			cr:   registryv1alpha1.DevfileRegistry{},
			want: DefaultNameOverride,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetNameOverride(&tt.cr)
			if result != tt.want {
				t.Errorf("func TestGetNameOverride(t *testing.T) {\n error: enablement value mismatch, expected: %v got: %v", tt.want, result)
			}
		})
	}
}

func TestGetFullnameOverride(t *testing.T) {
	tests := []struct {
		name string
		cr   registryv1alpha1.DevfileRegistry
		want string
	}{
		{
			name: "Case 1: Full app name override is set",
			cr: registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					FullnameOverride: "devfile-registry-test",
				},
			},
			want: "devfile-registry-test",
		},
		{
			name: "Case 2: Full app name override is unset",
			cr:   registryv1alpha1.DevfileRegistry{},
			want: DefaultFullnameOverride,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetFullnameOverride(&tt.cr)
			if result != tt.want {
				t.Errorf("func TestGetFullnameOverride(t *testing.T) {\n error: enablement value mismatch, expected: %v got: %v", tt.want, result)
			}
		})
	}
}

func Test_getAppName(t *testing.T) {
	tests := []struct {
		name string
		cr   *registryv1alpha1.DevfileRegistry
		want string
	}{
		{
			name: "Case 1: Default App Name",
			cr:   &registryv1alpha1.DevfileRegistry{},
			want: DefaultAppName,
		},
		{
			name: "Case 2: Overridden Short App Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "dr",
				},
			},
			want: "dr",
		},
		{
			name: "Case 3: Overridden Long App Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1-tf433",
				},
			},
			want: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1",
		},
		{
			name: "Case 4: CR set to nil",
			cr:   nil,
			want: DefaultAppName,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := getAppName(test.cr)
			if got != test.want {
				t.Errorf("\nGot: %v\nExpected: %v\n", got, test.want)
			}
		})
	}
}

func Test_getAppFullName(t *testing.T) {
	tests := []struct {
		name string
		cr   *registryv1alpha1.DevfileRegistry
		want string
	}{
		{
			name: "Case 1: Default App Full Name",
			cr:   &registryv1alpha1.DevfileRegistry{},
			want: DefaultAppName,
		},
		{
			name: "Case 2: Default App Full Name with overridden short App Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "dr",
				},
			},
			want: "dr",
		},
		{
			name: "Case 3: Default App Full Name with overridden long App Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1-tf433",
				},
			},
			want: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1",
		},
		{
			name: "Case 4: Overridden Short App Full Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					FullnameOverride: "dr",
				},
			},
			want: "dr",
		},
		{
			name: "Case 5: Overridden Long App Full Name",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					FullnameOverride: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1-tf433",
				},
			},
			want: "devfile-registry-testregistry-devfile-io-k8s-prow-environment1",
		},
		{
			name: "Case 6: Default App Full Name with short CR name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dr",
				},
			},
			want: fmt.Sprintf("%s-%s", "dr", DefaultAppName),
		},
		{
			name: "Case 7: Default App Full Name with long CR name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "testregistry-devfile-io-k8s-prow-af4325d2dcb2d0rte1-devfile-registry",
				},
			},
			want: "testregistry-devfile-io-k8s-prow-af4325d2dcb2d0rte1-devfile-reg",
		},
		{
			name: "Case 8: Default App Full Name with CR name contains default app name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "devfile-registry-test",
				},
			},
			want: "devfile-registry-test",
		},
		{
			name: "Case 9: Default App Full Name with CR name contains overridden app name",
			cr: &registryv1alpha1.DevfileRegistry{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dr-test",
				},
				Spec: registryv1alpha1.DevfileRegistrySpec{
					NameOverride: "dr",
				},
			},
			want: "dr-test",
		},
		{
			name: "Case 10: CR set to nil",
			cr:   nil,
			want: DefaultAppName,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := getAppFullName(test.cr)
			if got != test.want {
				t.Errorf("\nGot: %v\nExpected: %v\n", got, test.want)
			}
		})
	}
}

func TestIsIngressSkipped(t *testing.T) {

	tests := []struct {
		name string
		cr   *registryv1alpha1.DevfileRegistry
		want bool
	}{
		{
			name: "Case 1: Ingress skipped",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					K8s: registryv1alpha1.DevfileRegistrySpecK8sOnly{},
				},
			},
			want: true,
		},
		{
			name: "Case 2: Ingress set",
			cr: &registryv1alpha1.DevfileRegistry{
				Spec: registryv1alpha1.DevfileRegistrySpec{
					K8s: registryv1alpha1.DevfileRegistrySpecK8sOnly{
						IngressDomain: "test",
					},
				},
			},
			want: false,
		},
		{
			name: "Case 3: CR is nil",
			cr:   nil,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ingressSkipped := IsIngressSkipped(tt.cr)
			if ingressSkipped != tt.want {
				t.Errorf("TestIsIngressSkipped error: value mismatch, expected: %v got: %v", tt.want, ingressSkipped)
			}
		})
	}

}
