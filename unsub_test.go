//
// Copyright © 2011-2019 Guy M. Allard
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

package stompws

import (
	//"fmt"
	"log"

	//"os"
	"testing"
	//"time"
)

func TestUnSubNoHeader(t *testing.T) {
	n, _ = openConn(t)
	ch := login_headers
	ch = headersProtocol(ch, SPL_10) // To start
	conn, e = Connect(n, ch)
	if e != nil {
		t.Fatalf("TestUnSubNoHeader CONNECT Failed: e:<%q> connresponse:<%q>\n", e,
			conn.ConnectResponse)
	}
	//
	for ti, tv := range unsubNoHeaderDataList {
		conn.protocol = tv.proto // Cheat, fake all protocols
		e = conn.Unsubscribe(empty_headers)
		if e == nil {
			t.Fatalf("TestUnSubNoHeader[%d] proto:%s expected:%q got:nil\n",
				ti, sp, tv.exe)
		}
		if e != tv.exe {
			t.Fatalf("TestUnSubNoHeader[%d] proto:%s expected:%q got:%q\n",
				ti, sp, tv.exe, e)
		}
	}
	//
	checkReceived(t, conn, false)
	e = conn.Disconnect(empty_headers)
	checkDisconnectError(t, e)
	_ = closeConn(t, n)
	log.Printf("TestUnSubNoHeader %d tests complete.\n", len(subNoHeaderDataList))

}

func TestUnSubNoID(t *testing.T) {
	n, _ = openConn(t)
	ch := login_headers
	ch = headersProtocol(ch, SPL_10) // To start
	conn, e = Connect(n, ch)
	if e != nil {
		t.Fatalf("TestUnSubNoID CONNECT Failed: e:<%q> connresponse:<%q>\n", e,
			conn.ConnectResponse)
	}
	//
	for ti, tv := range unsubNoHeaderDataList {
		conn.protocol = tv.proto // Cheat, fake all protocols
		e = conn.Unsubscribe(empty_headers)
		if e != tv.exe {
			t.Fatalf("TestUnSubNoHeader[%d] proto:%s expected:%q got:%q\n",
				ti, sp, tv.exe, e)
		}
	}
	//
	checkReceived(t, conn, false)
	e = conn.Disconnect(empty_headers)
	checkDisconnectError(t, e)
	_ = closeConn(t, n)
	log.Printf("TestUnSubNoID %d tests complete.\n", len(unsubNoHeaderDataList))
}

func TestUnSubBool(t *testing.T) {
	n, _ = openConn(t)
	ch := login_headers
	ch = headersProtocol(ch, SPL_10) // To start
	conn, e = Connect(n, ch)
	if e != nil {
		t.Fatalf("CONNECT Failed: e:<%q> connresponse:<%q>\n", e,
			conn.ConnectResponse)
	}
	//
	for ti, tv := range unsubBoolDataList {
		conn.protoLock.Lock()
		conn.protocol = tv.proto // Cheat, fake all protocols
		conn.protoLock.Unlock()

		// SUBSCRIBE Phase (depending on test data)
		if tv.subfirst {
			// Do a real SUBSCRIBE
			// sc, e = conn.Subscribe
			sh := fixHeaderDest(tv.subh) // destination fixed if needed
			sc, e = conn.Subscribe(sh)
			if e == nil && sc == nil {
				t.Fatalf("TestUnSubBool[%d] SUBSCRIBE proto:%s expected OK, got <%v> <%v>\n",
					ti, conn.protocol, e, sc)
			}
			if sc == nil {
				t.Fatalf("TestUnSubBool[%d] SUBSCRIBE, proto:[%s], channel is nil\n",
					ti, tv.proto)
			}
			if e != tv.exe1 {
				t.Fatalf("TestUnSubBool[%d] SUBSCRIBE NEQCHECK proto:%s expected:%v got:%q\n",
					ti, tv.proto, tv.exe1, e)
			}
		}

		//fmt.Printf("fs,unsubh: <%v>\n", tv.unsubh)
		// UNSCRIBE Phase
		sh := fixHeaderDest(tv.unsubh) // destination fixed if needed
		e = conn.Unsubscribe(sh)
		if e != tv.exe2 {
			t.Fatalf("TestUnSubBool[%d] UNSUBSCRIBE NEQCHECK proto:%s expected:%v got:%q\n",
				ti, tv.proto, tv.exe2, e)
		}
	}
	//
	checkReceived(t, conn, true) // true for this test
	_ = conn.Disconnect(empty_headers)
	_ = closeConn(t, n)
	log.Printf("TestUnSubBool %d tests complete.\n", len(unsubBoolDataList))
}
