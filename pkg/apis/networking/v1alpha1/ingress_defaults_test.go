/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/util/intstr"

	"knative.dev/networking/pkg/apis/config"
)

var defaultMaxRevisionTimeout = time.Duration(config.DefaultMaxRevisionTimeoutSeconds) * time.Second

func TestIngressDefaulting(t *testing.T) {
	tests := []struct {
		name string
		in   *Ingress
		want *Ingress
	}{{
		name: "empty",
		in:   &Ingress{},
		want: &Ingress{
			Spec: IngressSpec{
				Visibility: IngressVisibilityExternalIP,
			},
		},
	}, {
		name: "has-visibility",
		in: &Ingress{
			Spec: IngressSpec{
				Visibility: IngressVisibilityClusterLocal,
			},
		},
		want: &Ingress{
			Spec: IngressSpec{
				Visibility: IngressVisibilityClusterLocal,
			},
		},
	}, {
		name: "split-timeout-defaulting",
		in: &Ingress{
			Spec: IngressSpec{
				Rules: []IngressRule{{
					HTTP: &HTTPIngressRuleValue{
						Paths: []HTTPIngressPath{{
							Splits: []IngressBackendSplit{{
								IngressBackend: IngressBackend{
									ServiceName:      "revision-000",
									ServiceNamespace: "default",
									ServicePort:      intstr.FromInt(8080),
								},
							}},
						}},
					},
				}},
				Visibility: IngressVisibilityExternalIP,
			},
		},
		want: &Ingress{
			Spec: IngressSpec{
				Rules: []IngressRule{{
					HTTP: &HTTPIngressRuleValue{
						Paths: []HTTPIngressPath{{
							Splits: []IngressBackendSplit{{
								IngressBackend: IngressBackend{
									ServiceName:      "revision-000",
									ServiceNamespace: "default",
									ServicePort:      intstr.FromInt(8080),
								},
								// Percent is filled in.
								Percent: 100,
							}},
						}},
					},
				}},
				Visibility: IngressVisibilityExternalIP,
			},
		},
	}, {
		name: "split-timeout-not-defaulting",
		in: &Ingress{
			Spec: IngressSpec{
				Rules: []IngressRule{{
					HTTP: &HTTPIngressRuleValue{
						Paths: []HTTPIngressPath{{
							Splits: []IngressBackendSplit{{
								IngressBackend: IngressBackend{
									ServiceName:      "revision-000",
									ServiceNamespace: "default",
									ServicePort:      intstr.FromInt(8080),
								},
								Percent: 30,
							}, {
								IngressBackend: IngressBackend{
									ServiceName:      "revision-001",
									ServiceNamespace: "default",
									ServicePort:      intstr.FromInt(8080),
								},
								Percent: 70,
							}},
						}},
					},
				}},
				Visibility: IngressVisibilityExternalIP,
			},
		},
		want: &Ingress{
			Spec: IngressSpec{
				Rules: []IngressRule{{
					HTTP: &HTTPIngressRuleValue{
						Paths: []HTTPIngressPath{{
							Splits: []IngressBackendSplit{{
								IngressBackend: IngressBackend{
									ServiceName:      "revision-000",
									ServiceNamespace: "default",
									ServicePort:      intstr.FromInt(8080),
								},
								// Percent is kept intact.
								Percent: 30,
							}, {
								IngressBackend: IngressBackend{
									ServiceName:      "revision-001",
									ServiceNamespace: "default",
									ServicePort:      intstr.FromInt(8080),
								},
								// Percent is kept intact.
								Percent: 70,
							}},
						}},
					},
				}},
				Visibility: IngressVisibilityExternalIP,
			},
		},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.in
			got.SetDefaults(context.Background())
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("SetDefaults (-want, +got) = %v", diff)
			}
		})
	}

}
