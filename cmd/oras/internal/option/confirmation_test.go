/*
Copyright The ORAS Authors.
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

package option

import (
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

func TestConfirmation_ApplyFlags(t *testing.T) {
	var test struct{ Confirmation }
	ApplyFlags(&test, pflag.NewFlagSet("oras-test", pflag.ExitOnError))
	if test.Confirmation.Confirmed != false {
		t.Fatalf("expecting Confirmed to be false but got: %v", test.Confirmation.Confirmed)
	}
}

func TestConfirmation_AskForConfirmation_ForciblyConfirmed(t *testing.T) {
	opts := Confirmation{
		Confirmed: true,
	}
	r := strings.NewReader("")

	got, err := opts.AskForConfirmation(r, "")
	if err != nil {
		t.Fatal("Confirmation.AskForConfirmation() error =", err)
	}
	if !reflect.DeepEqual(got, true) {
		t.Fatalf("Confirmation.AskForConfirmation() got %v, want %v", got, true)
	}
}

func TestConfirmation_AskForConfirmation_ManuallyConfirmed(t *testing.T) {
	opts := Confirmation{
		Confirmed: false,
	}

	r := strings.NewReader("yes")
	got, err := opts.AskForConfirmation(r, "")
	if err != nil {
		t.Fatal("Confirmation.AskForConfirmation() error =", err)
	}
	if !reflect.DeepEqual(got, true) {
		t.Fatalf("Confirmation.AskForConfirmation() got %v, want %v", got, true)
	}

	r = strings.NewReader("no")
	got, err = opts.AskForConfirmation(r, "")
	if err != nil {
		t.Fatal("Confirmation.AskForConfirmation() error =", err)
	}
	if !reflect.DeepEqual(got, false) {
		t.Fatalf("Confirmation.AskForConfirmation() got %v, want %v", got, false)
	}
}

func TestConfirmation_AskForConfirmation_ScanlnErr(t *testing.T) {
	r := strings.NewReader("yes no")
	opts := Confirmation{
		Confirmed: false,
	}

	_, err := opts.AskForConfirmation(r, "")
	expected := "expected newline"
	if err.Error() != expected {
		t.Fatalf("AskForConfirmation() error = %v, wantErr %v", err, expected)
	}
}
