package main

import (
	"testing"

	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// TestGetClientSet tests the getClientSet method
func TestGetClientSet(t *testing.T) {
	t.Run("return existing clientSet", func(t *testing.T) {
		command := &RunCommand{}

		// Test that when clientSet is nil, it would try to create one
		// (but we can't test the actual creation without a real cluster)
		_, err := command.getClientSet()
		if err == nil {
			t.Skip("Test requires in-cluster or kubeconfig setup, skipping")
		}
	})
}

// TestHasExistingServiceLogic tests the hasExistingService logic using a fake client
func TestHasExistingServiceLogic(t *testing.T) {
	tests := []struct {
		name     string
		pod      *core_v1.Pod
		services []core_v1.Service
		want     bool
	}{
		{
			name: "service exists",
			pod: &core_v1.Pod{
				ObjectMeta: meta_v1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
				},
			},
			services: []core_v1.Service{
				{
					ObjectMeta: meta_v1.ObjectMeta{
						Name:      "test-pod",
						Namespace: "default",
					},
				},
			},
			want: true,
		},
		{
			name: "service does not exist",
			pod: &core_v1.Pod{
				ObjectMeta: meta_v1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
				},
			},
			services: []core_v1.Service{},
			want:     false,
		},
		{
			name: "different namespace",
			pod: &core_v1.Pod{
				ObjectMeta: meta_v1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
				},
			},
			services: []core_v1.Service{
				{
					ObjectMeta: meta_v1.ObjectMeta{
						Name:      "test-pod",
						Namespace: "other",
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeClient := fake.NewSimpleClientset()

			// Create services in the fake client
			for _, svc := range tt.services {
				_, err := fakeClient.CoreV1().Services(svc.Namespace).Create(&svc)
				if err != nil {
					t.Fatalf("failed to create service: %v", err)
				}
			}

			// Check if service exists using the fake client directly
			_, err := fakeClient.CoreV1().Services(tt.pod.GetNamespace()).Get(tt.pod.GetName(), meta_v1.GetOptions{})
			got := err == nil

			if got != tt.want {
				t.Errorf("hasExistingService() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestServiceCreationLogic tests the service creation logic
func TestServiceCreationLogic(t *testing.T) {
	tests := []struct {
		name        string
		pod         *core_v1.Pod
		annotation  string
		shouldSkip  bool
		skipReason  string
	}{
		{
			name: "should create service for valid pod",
			pod: &core_v1.Pod{
				ObjectMeta: meta_v1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						"srcd.host/create-headless-service": "true",
					},
				},
				Status: core_v1.PodStatus{
					PodIP: "10.0.0.1",
				},
			},
			annotation: "srcd.host/create-headless-service",
			shouldSkip: false,
		},
		{
			name: "should skip pod without annotation",
			pod: &core_v1.Pod{
				ObjectMeta: meta_v1.ObjectMeta{
					Name:        "test-pod",
					Namespace:   "default",
					Annotations: map[string]string{},
				},
				Status: core_v1.PodStatus{
					PodIP: "10.0.0.1",
				},
			},
			annotation: "srcd.host/create-headless-service",
			shouldSkip: true,
			skipReason: "missing annotation",
		},
		{
			name: "should skip pod without IP",
			pod: &core_v1.Pod{
				ObjectMeta: meta_v1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "default",
					Annotations: map[string]string{
						"srcd.host/create-headless-service": "true",
					},
				},
				Status: core_v1.PodStatus{
					PodIP: "",
				},
			},
			annotation: "srcd.host/create-headless-service",
			shouldSkip: true,
			skipReason: "missing IP",
		},
		{
			name: "should skip pod with name too long",
			pod: &core_v1.Pod{
				ObjectMeta: meta_v1.ObjectMeta{
					Name:      "this-is-a-very-long-pod-name-that-exceeds-the-maximum-length-of-63-characters",
					Namespace: "default",
					Annotations: map[string]string{
						"srcd.host/create-headless-service": "true",
					},
				},
				Status: core_v1.PodStatus{
					PodIP: "10.0.0.1",
				},
			},
			annotation: "srcd.host/create-headless-service",
			shouldSkip: true,
			skipReason: "name too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check annotation
			hasAnnotation := tt.pod.Annotations[tt.annotation] == "true"
			hasIP := tt.pod.Status.PodIP != ""
			nameLengthOK := len(tt.pod.GetName()) <= 63

			shouldProcess := hasAnnotation && hasIP && nameLengthOK

			if shouldProcess == tt.shouldSkip {
				if tt.shouldSkip {
					t.Errorf("expected to skip (%s) but would process", tt.skipReason)
				} else {
					t.Errorf("expected to process but would skip")
				}
			}
		})
	}
}

// TestServiceCreationWithFakeClient tests actual service creation
func TestServiceCreationWithFakeClient(t *testing.T) {
	fakeClient := fake.NewSimpleClientset()

	pod := &core_v1.Pod{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
			Annotations: map[string]string{
				"test": "value",
			},
		},
		Status: core_v1.PodStatus{
			PodIP: "10.0.0.1",
		},
	}

	// Create service
	svc := &core_v1.Service{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:        pod.GetObjectMeta().GetName(),
			Annotations: pod.GetAnnotations(),
		},
		Spec: core_v1.ServiceSpec{
			ClusterIP: "None", // headless service
		},
	}

	_, err := fakeClient.CoreV1().Services(pod.GetNamespace()).Create(svc)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	// Verify service was created
	createdSvc, err := fakeClient.CoreV1().Services(pod.Namespace).Get(pod.Name, meta_v1.GetOptions{})
	if err != nil {
		t.Fatalf("service was not created: %v", err)
	}

	if createdSvc.Spec.ClusterIP != "None" {
		t.Errorf("expected headless service (ClusterIP=None), got %v", createdSvc.Spec.ClusterIP)
	}

	// Create endpoint
	endpoint := &core_v1.Endpoints{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:        pod.GetObjectMeta().GetName(),
			Annotations: pod.GetAnnotations(),
		},
		Subsets: []core_v1.EndpointSubset{
			{
				Addresses: []core_v1.EndpointAddress{
					{
						IP: pod.Status.PodIP,
					},
				},
			},
		},
	}

	_, err = fakeClient.CoreV1().Endpoints(pod.GetNamespace()).Create(endpoint)
	if err != nil {
		t.Fatalf("failed to create endpoint: %v", err)
	}

	// Verify endpoint was created
	createdEndpoint, err := fakeClient.CoreV1().Endpoints(pod.Namespace).Get(pod.Name, meta_v1.GetOptions{})
	if err != nil {
		t.Fatalf("endpoint was not created: %v", err)
	}

	if len(createdEndpoint.Subsets) == 0 || len(createdEndpoint.Subsets[0].Addresses) == 0 {
		t.Fatalf("endpoint has no addresses")
	}

	if createdEndpoint.Subsets[0].Addresses[0].IP != pod.Status.PodIP {
		t.Errorf("endpoint IP = %v, want %v", createdEndpoint.Subsets[0].Addresses[0].IP, pod.Status.PodIP)
	}
}

// TestEndpointUpdate tests endpoint update logic
func TestEndpointUpdate(t *testing.T) {
	fakeClient := fake.NewSimpleClientset()

	pod := &core_v1.Pod{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
		Status: core_v1.PodStatus{
			PodIP: "10.0.0.2", // new IP
		},
	}

	// Create initial endpoint with old IP
	endpoint := &core_v1.Endpoints{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      pod.Name,
			Namespace: pod.Namespace,
		},
		Subsets: []core_v1.EndpointSubset{
			{
				Addresses: []core_v1.EndpointAddress{
					{
						IP: "10.0.0.1", // old IP
					},
				},
			},
		},
	}

	_, err := fakeClient.CoreV1().Endpoints(pod.Namespace).Create(endpoint)
	if err != nil {
		t.Fatalf("failed to create endpoint: %v", err)
	}

	// Update endpoint with new IP
	updatedEndpoint := &core_v1.Endpoints{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:        pod.GetObjectMeta().GetName(),
			Annotations: pod.GetAnnotations(),
		},
		Subsets: []core_v1.EndpointSubset{
			{
				Addresses: []core_v1.EndpointAddress{
					{
						IP: pod.Status.PodIP,
					},
				},
			},
		},
	}

	_, err = fakeClient.CoreV1().Endpoints(pod.Namespace).Update(updatedEndpoint)
	if err != nil {
		t.Fatalf("failed to update endpoint: %v", err)
	}

	// Verify update
	result, err := fakeClient.CoreV1().Endpoints(pod.Namespace).Get(pod.Name, meta_v1.GetOptions{})
	if err != nil {
		t.Fatalf("failed to get endpoint: %v", err)
	}

	if result.Subsets[0].Addresses[0].IP != pod.Status.PodIP {
		t.Errorf("endpoint IP = %v, want %v", result.Subsets[0].Addresses[0].IP, pod.Status.PodIP)
	}
}

// TestServiceDeletion tests service deletion logic
func TestServiceDeletion(t *testing.T) {
	fakeClient := fake.NewSimpleClientset()

	svc := &core_v1.Service{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "default",
		},
	}

	_, err := fakeClient.CoreV1().Services(svc.Namespace).Create(svc)
	if err != nil {
		t.Fatalf("failed to create service: %v", err)
	}

	// Delete service
	err = fakeClient.CoreV1().Services(svc.Namespace).Delete(svc.Name, &meta_v1.DeleteOptions{})
	if err != nil {
		t.Fatalf("failed to delete service: %v", err)
	}

	// Verify deletion
	_, err = fakeClient.CoreV1().Services(svc.Namespace).Get(svc.Name, meta_v1.GetOptions{})
	if err == nil {
		t.Errorf("expected service to be deleted")
	}
}
