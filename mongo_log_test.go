// Copyright 2018 Kuei-chun Chen. All rights reserved.

package gox

import (
	"testing"
)

func TestGet(t *testing.T) {
	log := `2019-05-26T10:51:53.448-0400 I COMMAND  [conn23075] command keyhole.cars appName: "MongoDB Shell" command: find { find: "cars", filter: { color: "Red", style: "Truck" }, sort: { year: 1.0 }, lsid: { id: UUID("36f308ad-ab58-4053-bd0f-d2d815269f4e") }, $clusterTime: { clusterTime: Timestamp(1558882305, 1), signature: { hash: BinData(0, 0000000000000000000000000000000000000000), keyId: 0 } }, $db: "keyhole" } planSummary: IXSCAN { color: 1, style: 1, year: 1 } cursorid:66325545984 keysExamined:101 docsExamined:101 fromMultiPlanner:1 numYields:4 nreturned:101 reslen:34418 locks:{ Global: { acquireCount: { r: 5 } }, Database: { acquireCount: { r: 5 } }, Collection: { acquireCount: { r: 5 } } } storage:{ data: { bytesRead: 14745, timeReadingMicros: 5342 } } protocol:op_msg 7ms`
	filter := `{ color: "Red", style: "Truck" }`
	mlog := NewMongoLog(log)
	matched := mlog.Get("filter:")
	if matched != filter {
		t.Fatal("Expected", filter, "but got", matched)
	}

	log = `2019-05-26T10:51:53.448-0400 I COMMAND  [conn23075] command keyhole.cars appName: "MongoDB Shell" command: find { find: "cars", filter: { number: 100, value: 456 }, sort: { year: 1.0 }, lsid: { id: UUID("36f308ad-ab58-4053-bd0f-d2d815269f4e") }, $clusterTime: { clusterTime: Timestamp(1558882305, 1), signature: { hash: BinData(0, 0000000000000000000000000000000000000000), keyId: 0 } }, $db: "keyhole" } planSummary: IXSCAN { color: 1, style: 1, year: 1 } cursorid:66325545984 keysExamined:101 docsExamined:101 fromMultiPlanner:1 numYields:4 nreturned:101 reslen:34418 locks:{ Global: { acquireCount: { r: 5 } }, Database: { acquireCount: { r: 5 } }, Collection: { acquireCount: { r: 5 } } } storage:{ data: { bytesRead: 14745, timeReadingMicros: 5342 } } protocol:op_msg 7ms`
	filter = `{ number: 100, value: 456 }`
	mlog = NewMongoLog(log)
	matched = mlog.Get("filter:")
	if matched != filter {
		t.Fatal("Expected", filter, "but got", matched)
	}
}
