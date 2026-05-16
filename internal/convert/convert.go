package convert

import (
	"github.com/aga-absolut/Vault-System/internal/models"
	pb "github.com/aga-absolut/Vault-System/proto/vault_system"
)

// ToRecord converts a protobuf record into an internal model record.
func ToRecord(pbRecord *pb.Record) *models.Record {
	if pbRecord == nil {
		return nil
	}

	return &models.Record{
		ID:   pbRecord.Id,
		Type: pbRecord.Type,
		Meta: pbRecord.Meta,
		Data: pbRecord.Data,
	}
}

// ToProtoRecord converts an internal model record into a protobuf record.
func ToProtoRecord(record *models.Record) *pb.Record {
	if record == nil {
		return nil
	}

	return &pb.Record{
		Id:   record.ID,
		Type: record.Type,
		Meta: record.Meta,
		Data: record.Data,
	}
}

// ToListProtoRecord converts a slice of model records into protobuf records.
func ToListProtoRecord(list []models.Record) []*pb.Record {
	if list == nil {
		return nil
	}

	var sliceRecords []*pb.Record

	for _, record := range list {
		protoData := ToProtoRecord(&record)
		sliceRecords = append(sliceRecords, protoData)
	}

	return sliceRecords
}

// ToListRecord converts a slice of protobuf records into model records.
func ToListRecord(protoList []*pb.Record) []models.Record {
	if protoList == nil {
		return nil
	}

	var sliceProtoRecords []models.Record

	for _, record := range protoList {
		data := ToRecord(record)
		sliceProtoRecords = append(sliceProtoRecords, *data)
	}

	return sliceProtoRecords
}