//
// Copyright © 2011-2016 Guy M. Allard
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package stompngo

import (
	"os"
	"testing"
)

type verData struct {
	ch Headers // Client headers
	sh Headers // Server headers
	e  error   // Expected error
}

var verChecks = []verData{
	{Headers{HK_ACCEPT_VERSION, SPL_11}, Headers{HK_VERSION, SPL_11}, nil},
	{Headers{}, Headers{}, nil},
	{Headers{HK_ACCEPT_VERSION, "1.0,1.1,1.2"}, Headers{HK_VERSION, SPL_12}, nil},
	{Headers{HK_ACCEPT_VERSION, "1.3"}, Headers{HK_VERSION, "1.3"}, EBADVERSVR},
	{Headers{HK_ACCEPT_VERSION, "1.3"}, Headers{HK_VERSION, "1.1"}, EBADVERCLI},
	{Headers{HK_ACCEPT_VERSION, "1.0,1.1,1.2"}, Headers{}, nil},
}

/*
	ConnDisc Test: net.Conn.
*/
func TestConnDiscNetconn(t *testing.T) {
	n, _ := openConn(t)
	_ = closeConn(t, n)
}

/*
	ConnDisc Test: stompngo.Connect.
*/
func TestConnDiscStompConn(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, e := Connect(n, ch)
	if e != nil {
		t.Fatalf("Expected no connect error, got [%v]\n", e)
	}
	if conn == nil {
		t.Fatalf("Expected a connection, got [nil]\n")
	}
	if conn.ConnectResponse.Command != CONNECTED {
		t.Fatalf("Expected command [%v], got [%v]\n", CONNECTED,
			conn.ConnectResponse.Command)
	}
	if !conn.connected {
		t.Fatalf("Expected connected [true], got [false]\n")
	}
	if !conn.Connected() {
		t.Fatalf("Expected connected [true], got [false]\n")
	}
	//
	if conn.Session() == "" {
		t.Fatalf("Expected connected session, got [default value]\n")
	}
	//
	if conn.SendTickerInterval() != 0 {
		t.Fatalf("Expected zero SendTickerInterval, got [%v]\n",
			conn.SendTickerInterval())
	}
	if conn.ReceiveTickerInterval() != 0 {
		t.Fatalf("Expected zero ReceiveTickerInterval, got [%v]\n",
			conn.SendTickerInterval())
	}
	if conn.SendTickerCount() != 0 {
		t.Fatalf("Expected zero SendTickerCount, got [%v]\n",
			conn.SendTickerCount())
	}
	if conn.ReceiveTickerCount() != 0 {
		t.Fatalf("Expected zero ReceiveTickerCount, got [%v]\n",
			conn.SendTickerCount())
	}
	//
	if conn.FramesRead() != 1 {
		t.Fatalf("Expected 1 frame read, got [%d]\n", conn.FramesRead())
	}
	if conn.BytesRead() <= 0 {
		t.Fatalf("Expected non-zero bytes read, got [%d]\n", conn.BytesRead())
	}
	if conn.FramesWritten() != 1 {
		t.Fatalf("Expected 1 frame written, got [%d]\n", conn.FramesWritten())
	}
	if conn.BytesWritten() <= 0 {
		t.Fatalf("Expected non-zero bytes written, got [%d]\n",
			conn.BytesWritten())
	}
	if conn.Running().Nanoseconds() == 0 {
		t.Fatalf("Expected non-zero runtime, got [0]\n")
	}
	//
	_ = conn.Disconnect(empty_headers)
	if conn.Connected() {
		t.Fatalf("Expected connected [false], got [true]\n")
	}
	_ = closeConn(t, n)
}

/*
	ConnDisc Test: stompngo.Disconnect.
*/
func TestConnDiscStompDisc(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)
	e := conn.Disconnect(empty_headers)
	if e != nil {
		t.Fatalf("Expected no disconnect error, got [%v]\n", e)
	}
	_ = closeConn(t, n)
}

/*
	ConnDisc Test: stompngo.Disconnect with client bypassing a receipt.
*/
func TestConnDiscNoDiscReceipt(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)
	e := conn.Disconnect(NoDiscReceipt)
	if e != nil {
		t.Fatalf("Expected no disconnect error, got [%v]\n", e)
	}
	if conn.DisconnectReceipt.Message.Command != "" {
		t.Fatalf("Expected no disconnect receipt command, got [%v]\n",
			conn.DisconnectReceipt.Message.Command)
	}
	_ = closeConn(t, n)
}

/*
	ConnDisc Test: stompngo.Disconnect with receipt requested.
*/
func TestConnDiscStompDiscReceipt(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)
	rid := "my-receipt-001"
	e := conn.Disconnect(Headers{HK_RECEIPT, rid})
	if e != nil {

		t.Fatalf("Expected no disconnect error, got [%v]\n", e)
	}
	if conn.DisconnectReceipt.Error != nil {
		t.Fatalf("Expected no receipt error, got [%v]\n",
			conn.DisconnectReceipt.Error)
	}
	md := conn.DisconnectReceipt.Message
	irid, ok := md.Headers.Contains(HK_RECEIPT_ID)
	if !ok {
		t.Fatalf("Expected receipt-id, not received\n")
	}
	if rid != irid {
		t.Fatalf("Expected receipt-id [%q], got [%q]\n", rid, irid)
	}
	_ = closeConn(t, n)
}

/*
	ConnDisc Test: Body Length of CONNECTED response.
*/
func TestConnBodyLen(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)

	conn, e := Connect(n, ch)
	if e != nil {
		t.Fatalf("Expected no connect error, got [%v]\n", e)
	}
	if len(conn.ConnectResponse.Body) != 0 {
		t.Fatalf("Expected body length 0, got [%v]\n",
			len(conn.ConnectResponse.Body))
	}
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Conn11 Test: Test 1.1+ Connection.
*/
func TestConn11p(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, e := Connect(n, ch)
	if e != nil {
		t.Fatalf("Expected no connect error, got [%v]\n", e)
	}
	v := os.Getenv("STOMP_TEST11p")
	if v != "" {
		switch v {
		case SPL_12:
			if conn.Protocol() != SPL_12 {
				t.Fatalf("Expected protocol %v, got [%v]\n", SPL_12, conn.Protocol())
			}
		default:
			if conn.Protocol() != SPL_11 {
				t.Fatalf("Expected protocol %v, got [%v]\n", SPL_11, conn.Protocol())
			}
		}
	} else {
		if conn.Protocol() != SPL_10 {
			t.Fatalf("Expected protocol %v, got [%v]\n", SPL_10, conn.Protocol())
		}
	}
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	Conn11Receipt Test: Test receipt not allowed on connect.
*/
func TestConn11Receipt(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	ch = ch.Add(HK_RECEIPT, "abcd1234")
	_, e := Connect(n, ch)
	if e == nil {
		t.Fatalf("Expected connect error, got nil\n")
	}
	if e != ENORECPT {
		t.Fatalf("Expected [%v], got [%v]\n", ENORECPT, e)
	}
	_ = closeConn(t, n)
}

/*
	ConnDisc Test: ECONBAD
*/
func TestConnEconBad(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, e := Connect(n, ch)
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
	//
	e = conn.Abort(empty_headers)
	if e != ECONBAD {
		t.Fatalf("Abort expected [%v] got [%v]\n", ECONBAD, e)
	}
	e = conn.Ack(empty_headers)
	if e != ECONBAD {
		t.Fatalf("Ack expected [%v] got [%v]\n", ECONBAD, e)
	}
	e = conn.Begin(empty_headers)
	if e != ECONBAD {
		t.Fatalf("Begin expected [%v] got [%v]\n", ECONBAD, e)
	}
	e = conn.Commit(empty_headers)
	if e != ECONBAD {
		t.Fatalf("Commit expected [%v] got [%v]\n", ECONBAD, e)
	}
	e = conn.Nack(empty_headers)
	if e != ECONBAD {
		t.Fatalf("Nack expected [%v] got [%v]\n", ECONBAD, e)
	}
	e = conn.Send(empty_headers, "")
	if e != ECONBAD {
		t.Fatalf("Send expected [%v] got [%v]\n", ECONBAD, e)
	}
	_, e = conn.Subscribe(empty_headers)
	if e != ECONBAD {
		t.Fatalf("Subscribe expected [%v] got [%v]\n", ECONBAD, e)
	}
	e = conn.Unsubscribe(empty_headers)
	if e != ECONBAD {
		t.Fatalf("Unsubscribe expected [%v] got [%v]\n", ECONBAD, e)
	}
}

/*
	ConnDisc Test: ECONBAD
*/
func TestConnEconDiscDone(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, e := Connect(n, ch)
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
	//
	e = conn.Disconnect(empty_headers)
	if e != ECONBAD {
		t.Fatalf("Previous disconnect expected [%v] got [%v]\n", ECONBAD, e)
	}
}

/*
	ConnDisc Test: setProtocolLevel
*/
func TestConnSetProtocolLevel(t *testing.T) {
	n, _ := openConn(t)
	ch := check11(TEST_HEADERS)
	conn, _ := Connect(n, ch)
	//
	for i, v := range verChecks {
		conn.protocol = SPL_10 // reset
		e := conn.setProtocolLevel(v.ch, v.sh)
		if e != v.e {
			t.Fatalf("Verdata Item [%d], expected [%v], got [%v]\n", i, v.e, e)
		}
	}
	//
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
}

/*
	ConnDisc Test: connRespData
*/
func TestConnRespData(t *testing.T) {

	for i, f := range frames {
		_, e := connectResponse(f.data)
		if e != f.resp {
			t.Fatalf("Index [%v], expected [%v], got [%v]\n", i, f.resp, e)
		}
	}
}
