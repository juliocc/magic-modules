package compute

import (
	"testing"

	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
)

func TestComputeInstance_networkIPCustomizedDiff(t *testing.T) {
	t.Parallel()

	d := &tpgresource.ResourceDiffMock{
		Before: map[string]interface{}{
			"network_interface.#": 0,
		},
		After: map[string]interface{}{
			"network_interface.#": 1,
		},
	}

	err := forceNewIfNetworkIPNotUpdatableFunc(d)
	if err != nil {
		t.Error(err)
	}

	if d.IsForceNew {
		t.Errorf("Expected not force new if network_interface array size changes")
	}

	type NetworkInterface struct {
		Network           string
		Subnetwork        string
		SubnetworkProject string
		NetworkIP         string
	}
	NIBefore := NetworkInterface{
		Network:           "a",
		Subnetwork:        "a",
		SubnetworkProject: "a",
		NetworkIP:         "a",
	}

	cases := map[string]struct {
		ExpectedForceNew bool
		Before           NetworkInterface
		After            NetworkInterface
	}{
		"NetworkIP only change": {
			ExpectedForceNew: true,
			Before:           NIBefore,
			After: NetworkInterface{
				Network:           "a",
				Subnetwork:        "a",
				SubnetworkProject: "a",
				NetworkIP:         "b",
			},
		},
		"NetworkIP and Network change": {
			ExpectedForceNew: false,
			Before:           NIBefore,
			After: NetworkInterface{
				Network:           "b",
				Subnetwork:        "a",
				SubnetworkProject: "a",
				NetworkIP:         "b",
			},
		},
		"NetworkIP and Subnetwork change": {
			ExpectedForceNew: false,
			Before:           NIBefore,
			After: NetworkInterface{
				Network:           "a",
				Subnetwork:        "b",
				SubnetworkProject: "a",
				NetworkIP:         "b",
			},
		},
		"NetworkIP and SubnetworkProject change": {
			ExpectedForceNew: false,
			Before:           NIBefore,
			After: NetworkInterface{
				Network:           "a",
				Subnetwork:        "a",
				SubnetworkProject: "b",
				NetworkIP:         "b",
			},
		},
		"All change": {
			ExpectedForceNew: false,
			Before:           NIBefore,
			After: NetworkInterface{
				Network:           "b",
				Subnetwork:        "b",
				SubnetworkProject: "b",
				NetworkIP:         "b",
			},
		},
		"No change": {
			ExpectedForceNew: false,
			Before:           NIBefore,
			After: NetworkInterface{
				Network:           "a",
				Subnetwork:        "a",
				SubnetworkProject: "a",
				NetworkIP:         "a",
			},
		},
	}

	for tn, tc := range cases {
		d := &tpgresource.ResourceDiffMock{
			Before: map[string]interface{}{
				"network_interface.#":                    1,
				"network_interface.0.network":            tc.Before.Network,
				"network_interface.0.subnetwork":         tc.Before.Subnetwork,
				"network_interface.0.subnetwork_project": tc.Before.SubnetworkProject,
				"network_interface.0.network_ip":         tc.Before.NetworkIP,
			},
			After: map[string]interface{}{
				"network_interface.#":                    1,
				"network_interface.0.network":            tc.After.Network,
				"network_interface.0.subnetwork":         tc.After.Subnetwork,
				"network_interface.0.subnetwork_project": tc.After.SubnetworkProject,
				"network_interface.0.network_ip":         tc.After.NetworkIP,
			},
		}
		err := forceNewIfNetworkIPNotUpdatableFunc(d)
		if err != nil {
			t.Error(err)
		}
		if tc.ExpectedForceNew != d.IsForceNew {
			t.Errorf("%v: expected d.IsForceNew to be %v, but was %v", tn, tc.ExpectedForceNew, d.IsForceNew)
		}
	}
}

func TestFlattenReplicaZones(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		Input    []string
		Expected []string
	}{
		"full self-links": {
			Input: []string{
				"https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-a",
				"https://www.googleapis.com/compute/v1/projects/my-project/zones/us-central1-b",
			},
			Expected: []string{"us-central1-a", "us-central1-b"},
		},
		"partial self-links": {
			Input: []string{
				"projects/my-project/zones/us-central1-a",
				"projects/my-project/zones/us-central1-b",
			},
			Expected: []string{"us-central1-a", "us-central1-b"},
		},
		"short names": {
			Input:    []string{"us-central1-a", "us-central1-b"},
			Expected: []string{"us-central1-a", "us-central1-b"},
		},
		"empty": {
			Input:    []string{},
			Expected: []string{},
		},
	}

	for tn, tc := range cases {
		actualInterfaces := flattenReplicaZones(tc.Input)
		
		if len(actualInterfaces) != len(tc.Expected) {
			t.Errorf("%s: expected length %d, got %d", tn, len(tc.Expected), len(actualInterfaces))
			continue
		}

		for i, expectedVal := range tc.Expected {
			actualVal, ok := actualInterfaces[i].(string)
			if !ok {
				t.Errorf("%s: element %d is not a string", tn, i)
				continue
			}
			if actualVal != expectedVal {
				t.Errorf("%s: at index %d expected %q, got %q", tn, i, expectedVal, actualVal)
			}
		}
	}
}
