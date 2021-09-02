package main

import (
	"context"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func pod(labels map[string]string, namespace, image string, phase v1.PodPhase, date metav1.Time) *v1.Pod {
	return &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Namespace: namespace, Labels: labels, Name: image, CreationTimestamp: date},
		Spec:       v1.PodSpec{Containers: []v1.Container{{Image: image}}}, Status: v1.PodStatus{Phase: phase},
	}
}

func TestReap(t *testing.T) {
	labelsone := map[string]string{"workflows.argoproj.io/completed": "true"}
	labelstwo := map[string]string{"workflows.argoproj.io/completed": "false"}
	var tests = []struct {
		description string
		objs        []runtime.Object
	}{
		{"600+ days old failed", []runtime.Object{pod(labelstwo, "asd", "a", v1.PodFailed, metav1.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)), pod(labelstwo, "asdasd", "b", v1.PodFailed, metav1.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC))}},
		{"600+ days old succeeded", []runtime.Object{pod(labelsone, "asd", "a", v1.PodSucceeded, metav1.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)), pod(labelsone, "asdasd", "b", v1.PodSucceeded, metav1.Date(2017, 1, 1, 0, 0, 0, 0, time.UTC))}},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			client := fake.NewSimpleClientset(test.objs...)

			p, _ := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

			if len(p.Items) == 0 {
				t.Errorf("Expected 2 pods, got %d", len(p.Items))
			}

			Reap(600, 600, []string{"asd", "asdasd"}, client.CoreV1(), false)

			p, _ = client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

			if len(p.Items) != 0 {
				t.Errorf("Expected 0 pods, got %d", len(p.Items))
			}
		})

		t.Run(test.description, func(t *testing.T) {
			client := fake.NewSimpleClientset(test.objs...)

			p, _ := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

			if len(p.Items) == 0 {
				t.Errorf("Expected 2 pods, got %d", len(p.Items))
			}

			Reap(600, 600, []string{"asd", "asdasd"}, client.CoreV1(), true)

			p, _ = client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

			if len(p.Items) == 0 {
				t.Errorf("Expected 2 pods, got %d", len(p.Items))
			}
		})

		t.Run(test.description, func(t *testing.T) {
			client := fake.NewSimpleClientset(test.objs...)

			p, _ := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

			if len(p.Items) != 2 {
				t.Errorf("Expected 2 pods, got %d", len(p.Items))
			}

			Reap(6000000, 6000000, []string{"asd", "asdasd"}, client.CoreV1(), false)

			p, _ = client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

			if len(p.Items) == 0 {
				t.Errorf("Expected 2 pods, got %d", len(p.Items))
			}
		})

		t.Run(test.description, func(t *testing.T) {
			client := fake.NewSimpleClientset(test.objs...)

			p, _ := client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

			if len(p.Items) != 2 {
				t.Errorf("Expected 2 pods, got %d", len(p.Items))
			}

			Reap(600, 600, []string{}, client.CoreV1(), false)

			p, _ = client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

			if len(p.Items) == 0 {
				t.Errorf("Expected 2 pods, got %d", len(p.Items))
			}
		})

	}
}
