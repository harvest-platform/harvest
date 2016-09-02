package json

import (
	"encoding/json"
	"testing"
)

func TestIsInteger(t *testing.T) {
	tests := []struct {
		V float64
		T bool
	}{
		{0, true},
		{1, true},
		{-1, true},
		{-1.0, true},
		{1.0, true},
		{0.1, false},
		{1.0000001, false},
		{0.9999999, false},
	}

	for _, x := range tests {
		r := isInteger(x.V)

		if r && !x.T {
			t.Errorf("%f should not be an integer", x.V)
		} else if !r && x.T {
			t.Errorf("%f should be an integer", x.V)
		}
	}
}

var jsonRecord = []byte(`
{"dob":"2005-01-29","audiograms":[{"audiometry":{"transducers":["ER-3A Inserts"],"test_methods":["Conventional"],"specfreq_transducers":["Sennheiser HDA200"],"visit_date":"2010-09-14","audiometer":"GSI 61","reliability":"Good"},"visit_date":"2010-09-14","puretone":{"ptas":[{"ear":"Left","conditions":["Air"],"value":50},{"ear":"Left","conditions":["Bone"],"value":33},{"ear":"Right","conditions":["Air"],"value":8}],"results":[{"conditions":["Air","Masked"],"ear":"Left","freqs":{"750":{"masking_level":null,"response":50},"6000":{"masking_level":null,"response":60},"4000":{"masking_level":null,"response":55},"3000":{"masking_level":null,"response":45},"250":{"masking_level":null,"response":80},"8000":{"masking_level":null,"response":60},"500":{"masking_level":null,"response":65}}},{"conditions":["Air","Unmasked"],"ear":"Left","freqs":{"4000":{"masking_level":null,"response":50},"2000":{"masking_level":null,"response":40},"1500":{"masking_level":null,"response":45},"1000":{"masking_level":null,"response":45}}},{"conditions":["Bone","Masked"],"ear":"Left","freqs":{"4000":{"masking_level":null,"response":40},"2000":{"masking_level":null,"response":25},"500":{"masking_level":null,"response":35},"1000":{"masking_level":null,"response":40}}},{"conditions":["Air","Unmasked"],"ear":"Right","freqs":{"4000":{"masking_level":null,"response":15},"2000":{"masking_level":null,"response":10},"250":{"masking_level":null,"response":15},"8000":{"masking_level":null,"response":5},"500":{"masking_level":null,"response":5},"1000":{"masking_level":null,"response":10}}}]}}],"sex":"Male","alias":"1630061","race":"White","organization":"Vanderbilt","ethnicity":"Missing"}
`)

func TestUnmarshalInfer(t *testing.T) {
	_, err := UnmarshalInfer(jsonRecord)
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkInfer(b *testing.B) {
	var m map[string]interface{}
	json.Unmarshal(jsonRecord, &m)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Infer(m)
	}
}
