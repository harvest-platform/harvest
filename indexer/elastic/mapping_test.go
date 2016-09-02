package elastic

import (
	"testing"

	"github.com/harvest-platform/harvest/indexer/json"
)

var jsonRecord = []byte(`
{"dob":"2005-01-29","audiograms":[{"audiometry":{"transducers":["ER-3A Inserts"],"test_methods":["Conventional"],"specfreq_transducers":["Sennheiser HDA200"],"visit_date":"2010-09-14","audiometer":"GSI 61","reliability":"Good"},"visit_date":"2010-09-14","puretone":{"ptas":[{"ear":"Left","conditions":["Air"],"value":50},{"ear":"Left","conditions":["Bone"],"value":33},{"ear":"Right","conditions":["Air"],"value":8}],"results":[{"conditions":["Air","Masked"],"ear":"Left","freqs":{"750":{"masking_level":null,"response":50},"6000":{"masking_level":null,"response":60},"4000":{"masking_level":null,"response":55},"3000":{"masking_level":null,"response":45},"250":{"masking_level":null,"response":80},"8000":{"masking_level":null,"response":60},"500":{"masking_level":null,"response":65}}},{"conditions":["Air","Unmasked"],"ear":"Left","freqs":{"4000":{"masking_level":null,"response":50},"2000":{"masking_level":null,"response":40},"1500":{"masking_level":null,"response":45},"1000":{"masking_level":null,"response":45}}},{"conditions":["Bone","Masked"],"ear":"Left","freqs":{"4000":{"masking_level":null,"response":40},"2000":{"masking_level":null,"response":25},"500":{"masking_level":null,"response":35},"1000":{"masking_level":null,"response":40}}},{"conditions":["Air","Unmasked"],"ear":"Right","freqs":{"4000":{"masking_level":null,"response":15},"2000":{"masking_level":null,"response":10},"250":{"masking_level":null,"response":15},"8000":{"masking_level":null,"response":5},"500":{"masking_level":null,"response":5},"1000":{"masking_level":null,"response":10}}}]}}],"sex":"Male","alias":"1630061","race":"White","organization":"Vanderbilt","ethnicity":"Missing"}
`)

func TestGenerate(t *testing.T) {
	fields, err := json.UnmarshalInfer(jsonRecord)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("generate", func(t *testing.T) {
		Generate(fields)
	})
}

func BenchmarkGenerate(b *testing.B) {
	fields, err := json.UnmarshalInfer(jsonRecord)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Generate(fields)
	}
}
