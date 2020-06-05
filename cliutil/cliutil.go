// Copyright 2015 Osipov Konstantin <k.osipov.msk@gmail.com>. All rights reserved.
// license that can be found in the LICENSE file.

// This file is part of the application source code leveldb-cli
// This software provides a console interface to leveldb.

package cliutil

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math"

	"github.com/TomiHiltunen/geohash-golang"
	"gopkg.in/mgo.v2/bson"
)

const (
	// KeyTypeDevice <int32 folder ID> <int32 device ID> <file name> = FileInfo
	KeyTypeDevice = 0

	// KeyTypeGlobal <int32 folder ID> <file name> = VersionList
	KeyTypeGlobal = 1

	// KeyTypeBlock <int32 folder ID> <32 bytes hash> <§file name> = int32 (block index)
	KeyTypeBlock = 2

	// KeyTypeDeviceStatistic <device ID as string> <some string> = some value
	KeyTypeDeviceStatistic = 3

	// KeyTypeFolderStatistic <folder ID as string> <some string> = some value
	KeyTypeFolderStatistic = 4

	// KeyTypeVirtualMtime <int32 folder ID> <file name> = dbMtime
	KeyTypeVirtualMtime = 5

	// KeyTypeFolderIdx <int32 id> = string value
	KeyTypeFolderIdx = 6

	// KeyTypeDeviceIdx <int32 id> = string value
	KeyTypeDeviceIdx = 7

	// KeyTypeIndexID <int32 device ID> <int32 folder ID> = protocol.IndexID
	KeyTypeIndexID = 8

	// KeyTypeFolderMeta <int32 folder ID> = CountsSet
	KeyTypeFolderMeta = 9

	// KeyTypeMiscData <some string> = some value
	KeyTypeMiscData = 10

	// KeyTypeSequence <int32 folder ID> <int64 sequence number> = KeyTypeDevice key
	KeyTypeSequence = 11

	// KeyTypeNeed <int32 folder ID> <file name> = <nothing>
	KeyTypeNeed = 12

	// KeyTypeBlockList <block list hash> = BlockList
	KeyTypeBlockList = 13
)

// Converts data to a string
func ToString(format string, value []byte) string {
	switch format {
	case "ignore":
		return ""
	case "bson":
		return bsonToString(value)
	case "base64":
		return ToBase64(value)
	case "geohash":
		return geohashToString(value)
	case "int64":
		return int64ToString(value)
	case "float64":
		return float64ToString(value)
	case "synckey":
		return ToSyncthingKey(value)
	case "raw":
	default:
	}

	return string(value)
}

func nulString(bs []byte) string {
	for i := range bs {
		if bs[i] == 0 {
			return string(bs[:i])
		}
	}
	return string(bs)
}

func ToSyncthingKey(key []byte) string {
	if len(key) > 0 {

		switch key[0] {
		case KeyTypeDevice:
			folder := binary.BigEndian.Uint32(key[1:])
			device := binary.BigEndian.Uint32(key[1+4:])
			name := nulString(key[1+4+4:])
			return fmt.Sprintf("[device] F:%d D:%d N:%q", folder, device, name)
		case KeyTypeGlobal:
			folder := binary.BigEndian.Uint32(key[1:])
			name := nulString(key[1+4:])
			return fmt.Sprintf("[global] F:%d N:%q", folder, name)

		case KeyTypeBlock:
			folder := binary.BigEndian.Uint32(key[1:])
			hash := key[1+4 : 1+4+32]
			name := nulString(key[1+4+32:])
			return fmt.Sprintf("[block] F:%d H:%x N:%q", folder, hash, name)

		case KeyTypeDeviceStatistic:
			return fmt.Sprintf("[dstat] K:%x", key)

		case KeyTypeFolderStatistic:
			return fmt.Sprintf("[fstat] K:%x", key)

		case KeyTypeVirtualMtime:
			folder := binary.BigEndian.Uint32(key[1:])
			name := nulString(key[1+4:])
			return fmt.Sprintf("[mtime] F:%d N:%q", folder, name)

		case KeyTypeFolderIdx:
			key := binary.BigEndian.Uint32(key[1:])
			return fmt.Sprintf("[folderidx] K:%d", key)

		case KeyTypeDeviceIdx:
			key := binary.BigEndian.Uint32(key[1:])

			return fmt.Sprintf("[deviceidx] K:%d", key)

		case KeyTypeIndexID:
			device := binary.BigEndian.Uint32(key[1:])
			folder := binary.BigEndian.Uint32(key[5:])
			return fmt.Sprintf("[indexid] D:%d F:%d", device, folder)

		case KeyTypeFolderMeta:
			folder := binary.BigEndian.Uint32(key[1:])
			return fmt.Sprintf("[foldermeta] F:%d", folder)

		case KeyTypeMiscData:
			return fmt.Sprintf("[miscdata] K:%q", key[1:])

		case KeyTypeSequence:
			folder := binary.BigEndian.Uint32(key[1:])
			seq := binary.BigEndian.Uint64(key[5:])
			return fmt.Sprintf("[sequence] F:%d S:%d", folder, seq)

		case KeyTypeNeed:
			folder := binary.BigEndian.Uint32(key[1:])
			file := string(key[5:])
			return fmt.Sprintf("[need] F:%d V:%q", folder, file)

		case KeyTypeBlockList:
			return fmt.Sprintf("[blocklist] H:%x", key[1:])

		default:
			return fmt.Sprintf("[???]\t  %x", key)
		}
	}

	return string(key)
}

// Converts data from bson type to a string
func bsonToString(value []byte) string {
	var dst interface{}
	err := bson.Unmarshal(value, &dst)
	if err != nil {
		return "Error converting!"
	}

	return fmt.Sprintf("%#v", dst)
}

// Converts data from bytes to a base64 string
func ToBase64(value []byte) string {
	return fmt.Sprintf("%s", base64.StdEncoding.EncodeToString(value))
}

// Converts data from geohash type to a string
func geohashToString(value []byte) string {
	position := geohash.Decode(string(value))

	return fmt.Sprintf("lat: %f lng: %f", position.Center().Lat(), position.Center().Lng())
}

// Converts data from int64 type to a string
func int64ToString(value []byte) string {
	return fmt.Sprintf("%d", binary.BigEndian.Uint64(value))
}

// Converts data from float64 type to a string
func float64ToString(value []byte) string {
	return fmt.Sprintf("%f", math.Float64frombits(
		binary.LittleEndian.Uint64(value),
	))
}
