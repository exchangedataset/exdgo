package exdgo

import (
	"bytes"
	"testing"
	"time"
)

func prepareReplayRequest(t *testing.T) *ReplayRequest {
	cli := ClientParam{
		APIKey: "demo",
	}
	start, serr := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	if serr != nil {
		t.Fatalf("testing error: %v", serr)
	}
	end, serr := time.Parse(time.RFC3339, "2020-01-01T00:04:50Z")
	if serr != nil {
		t.Fatalf("testing error: %v", serr)
	}
	reqp := ReplayRequestParam{
		Filter: map[string][]string{
			"bitmex":   []string{"orderBookL2_XBTUSD"},
			"bitfinex": []string{"trades_tBTCUSD"},
		},
		Start: start,
		End:   end,
	}
	req, serr := Replay(cli, reqp)
	if serr != nil {
		t.Fatal(serr)
	}
	return req
}

func TestReplayDownloadAndStream(t *testing.T) {
	req := prepareReplayRequest(t)

	lines, serr := req.Download()
	if serr != nil {
		t.Fatal(serr)
	}
	if len(lines) == 0 {
		t.Fatal("lines len 0")
	}
	itr, serr := req.Stream()
	if serr != nil {
		t.Fatal(serr)
	}
	defer itr.Close()
	i := 0
	for {
		line, ok, serr := itr.Next()
		if !ok {
			if serr != nil {
				t.Fatal(serr)
			}
			break
		}
		if *line.Channel != *lines[i].Channel {
			t.Fatal("channel differ")
		}
		if line.Exchange != lines[i].Exchange {
			t.Fatal("exchange differ")
		}
		switch line.Message.(type) {
		case []byte:
			if bytes.Compare(line.Message.([]byte), lines[i].Message.([]byte)) != 0 {
				t.Fatal("message differ")
			}
		case map[string]interface{}:
			for name, val := range line.Message.(map[string]interface{}) {
				va, ok := lines[i].Message.(map[string]interface{})[name]
				if !ok {
					t.Fatal("message differ: key does not exist")
				}
				if va != val {
					t.Fatal("message differ: value differ")
				}
			}
		default:
			t.Fatal("default should not be called")
		}
		if line.Timestamp != lines[i].Timestamp {
			t.Fatal("timestamp differ")
		}
		if line.Type != lines[i].Type {
			t.Fatal("type differ")
		}
		i++
	}
	if len(lines) != i {
		t.Fatal("len(lines) != i")
	}
}
