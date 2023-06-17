package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/mamalmaleki/go-movie/gen"
	"github.com/mamalmaleki/go-movie/metadata/pkg/model"
	"google.golang.org/protobuf/proto"
)

var metadataList = make([]model.Metadata, 0)
var genMetadataList = make([]gen.Metadata, 0)

func createMetadataList() {
	for i := 0; i < 50; i++ {
		m := model.Metadata{
			ID:          "XYZ-" + fmt.Sprint(i),
			Title:       "The Movie 2",
			Description: "Sequel of the legendary The Movie",
			Director:    "Foo Bars",
		}
		metadataList = append(metadataList, m)
	}
}

func createGenMetadataList() {
	for i := 0; i < 50; i++ {
		m := gen.Metadata{
			Id:          "XYZ-" + fmt.Sprint(i),
			Title:       "The Movie 2",
			Description: "Sequel of the legendary The Movie",
			Director:    "Foo Bars",
		}
		genMetadataList = append(genMetadataList, m)
	}
}

var metadata = &model.Metadata{
	ID:          "123",
	Title:       "The Movie 2",
	Description: "Sequel of the legendary The Movie",
	Director:    "Foo Bars",
}

var genMetadata = &gen.Metadata{
	Id:          "123",
	Title:       "The Movie 2",
	Description: "Sequel of the legendary The Movie",
	Director:    "Foo Bars",
}

func main() {
	jsonBytes, err := serializeToJSON(metadata)
	if err != nil {
		panic(err)
	}

	xmlBytes, err := serializeToXML(metadata)
	if err != nil {
		panic(err)
	}

	protoBytes, err := serializeToProto(genMetadata)
	if err != nil {
		panic(err)
	}

	fmt.Printf("JSON size:\t%dB\n", len(jsonBytes))
	fmt.Printf("XML size:\t%dB\n", len(xmlBytes))
	fmt.Printf("Proto size:\t%dB\n", len(protoBytes))

	createMetadataList()
	jsonArrBytes, err := serializeToJSONArray(metadataList)
	if err != nil {
		panic(err)
	}

	xmlArrBytes, err := serializeToXMLArray(metadataList)
	if err != nil {
		panic(err)
	}

	createGenMetadataList()
	protoArrBytesLen, err := getLenSerializedProtoArray(genMetadataList)
	if err != nil {
		panic(err)
	}

	fmt.Printf("JSON Array size:\t%dB\n", len(jsonArrBytes))
	fmt.Printf("XML Array size:\t%dB\n", len(xmlArrBytes))
	fmt.Printf("Proto Array size:\t%dB\n", protoArrBytesLen)
}

func getLenSerializedProtoArray(list []gen.Metadata) (int, error) {
	sum := 0
	for _, g := range list {
		b, err := serializeToProto(&g)
		if err != nil {
			return 0, err
		}
		sum += len(b)
	}
	return sum, nil
}

func serializeToXMLArray(list []model.Metadata) ([]byte, error) {
	return xml.Marshal(list)
}

func serializeToJSONArray(list []model.Metadata) ([]byte, error) {
	return json.Marshal(list)
}

func serializeToProto(m *gen.Metadata) ([]byte, error) {
	return proto.Marshal(m)
}

func serializeToXML(m *model.Metadata) ([]byte, error) {
	return xml.Marshal(m)
}

func serializeToJSON(m *model.Metadata) ([]byte, error) {
	return json.Marshal(m)
}
