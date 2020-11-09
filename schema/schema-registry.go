package schema

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"

	"github.com/riferrei/srclient"
)

type SchemaRegistry interface {
	GetLastestSchemaVersion(topic string) (sc *Schema, err error)
	CreateSchema(topic, schema, format string, isKey bool) (sc *srclient.Schema, err error)
	CreateSchemaFromFile(topic, schemaFile, format string, isKey bool) (*Schema, error)
	GetSchemaFromMsg(msg []byte) (*Schema, error)
}

type schemaRegistry struct {
	client *srclient.SchemaRegistryClient
}

type Schema struct {
	*srclient.Schema
}

func NewSchema(address string) SchemaRegistry {
	// 2) Fetch the latest version of the schema, or create a new one if it is the first
	schemaRegistryClient := srclient.CreateSchemaRegistryClient(address)
	return &schemaRegistry{
		client: schemaRegistryClient,
	}
}

func (s *schemaRegistry) GetLastestSchemaVersion(topic string) (sc *Schema, err error) {
	var schema *srclient.Schema
	schema, err = s.client.GetLatestSchema(topic, false)
	if err != nil {
		return nil, err
	}
	return &Schema{schema}, nil
}

func (s *schemaRegistry) CreateSchema(topic, schema, format string, isKey bool) (sc *srclient.Schema, err error) {
	return s.client.CreateSchema(topic, schema, srclient.SchemaType(format), isKey)
}

func (s *schemaRegistry) CreateSchemaFromFile(topic, schemaFile, format string, isKey bool) (*Schema, error) {
	fmt.Println("voo")
	schemaBytes, err := ioutil.ReadFile(schemaFile)
	fmt.Printf("schemaBytes %v\n", string(schemaBytes))
	if err != nil {
		panic(fmt.Sprintf("Error read schema file %s", err))
	}
	sc, err := s.client.CreateSchema(topic, string(schemaBytes), srclient.SchemaType(format), isKey)
	if err != nil {
		return nil, err
	}
	return &Schema{sc}, nil
}

func (s *schemaRegistry) GetSchemaFromMsg(msg []byte) (*Schema, error) {
	schemaID := binary.BigEndian.Uint32(msg[1:5])
	sc, err := s.client.GetSchema(int(schemaID))
	if err != nil {
		return nil, fmt.Errorf("Error getting the schema with id '%d' %s", schemaID, err)
	}
	return &Schema{sc}, nil
}

func (s *Schema) GetSchemaID() []byte {
	schemaIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(schemaIDBytes, uint32(s.ID()))
	return schemaIDBytes
}

func (s *Schema) SerializeData(payload []byte) []byte {
	native, _, _ := s.Codec().NativeFromTextual(payload)
	valueBytes, _ := s.Codec().BinaryFromNative(nil, native)

	var recordValue []byte
	recordValue = append(recordValue, byte(0))
	recordValue = append(recordValue, s.GetSchemaID()...) // 4 bytes of schema id
	recordValue = append(recordValue, valueBytes...)

	return recordValue
}

func (s *Schema) DeserializeData(msg []byte) []byte {
	native, _, _ := s.Codec().NativeFromBinary(msg[5:])
	value, _ := s.Codec().TextualFromNative(nil, native)
	fmt.Printf("Here is the message %s\n", string(value))
	return value
}
