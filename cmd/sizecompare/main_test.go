package main

import "testing"

func BenchmarkSerializeToJson(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = serializeToJSON(metadata)
	}
}

func BenchmarkSerializeToXML(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = serializeToXML(metadata)
	}
}

func BenchmarkSerializeToProto(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = serializeToProto(genMetadata)
	}
}
