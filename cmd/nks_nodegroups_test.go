package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// The create body needs image_id and the labels availability_zone,
// boot_volume_type, boot_volume_size and ca_enable, or the API answers 400.
func TestCreateNodeGroupSendsRequiredImageAndLabels(t *testing.T) {
	c := &cobra.Command{}
	c.Flags().AddFlagSet(nksCreateNodeGroupCmd.Flags())
	for _, f := range []string{"image-id", "availability-zone", "boot-volume-size"} {
		fl := c.Flags().Lookup(f)
		if fl == nil {
			t.Fatalf("create-node-group has no --%s flag", f)
		}
		if _, ok := fl.Annotations[cobra.BashCompOneRequiredFlag]; !ok {
			t.Errorf("--%s is not required", f)
		}
	}
	defer func() {
		for _, f := range []string{"image-id", "availability-zone", "boot-volume-size", "label"} {
			fl := c.Flags().Lookup(f)
			fl.Value.Set(fl.DefValue)
			fl.Changed = false
		}
	}()
	c.Flags().Set("image-id", "img-1")
	c.Flags().Set("availability-zone", "az-a")
	c.Flags().Set("boot-volume-size", "30")
	c.Flags().Set("label", "user_script_v2=echo,ca_enable=True")

	b, _ := json.Marshal(nodeGroupInput(c, "ng", "flv-1", 1))
	got := string(b)
	for _, want := range []string{`"image_id":"img-1"`, `"flavor_id":"flv-1"`, `"availability_zone":"az-a"`,
		`"boot_volume_type":"General HDD"`, `"boot_volume_size":"30"`, `"ca_enable":"True"`, `"user_script_v2":"echo"`} {
		if !strings.Contains(got, want) {
			t.Errorf("body %s lacks %s", got, want)
		}
	}
}
