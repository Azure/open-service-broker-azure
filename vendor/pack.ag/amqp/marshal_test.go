package amqp

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

var exampleFrames = []struct {
	label string
	frame frame
}{
	{
		label: "transfer",
		frame: frame{
			type_:   frameTypeAMQP,
			channel: 10,
			body: &performTransfer{
				Handle:             34983,
				DeliveryID:         uint32Ptr(564),
				DeliveryTag:        []byte("foo tag"),
				MessageFormat:      uint32Ptr(34),
				Settled:            true,
				More:               true,
				ReceiverSettleMode: rcvSettle(ModeSecond),
				State:              &stateReceived{},
				Resume:             true,
				Aborted:            true,
				Batchable:          true,
				Payload:            []byte("very important payload"),
			},
		},
	},
}

func TestFrameMarshalUnmarshal(t *testing.T) {
	for _, tt := range exampleFrames {
		t.Run(tt.label, func(t *testing.T) {
			var buf buffer

			err := writeFrame(&buf, tt.frame)
			if err != nil {
				t.Fatalf("%+v", err)
			}

			header, err := parseFrameHeader(&buf)
			if err != nil {
				t.Fatalf("%+v", err)
			}

			want := tt.frame
			if header.Channel != want.channel {
				t.Errorf("Expected channel to be %d, but it is %d", want.channel, header.Channel)
			}
			if header.FrameType != want.type_ {
				t.Errorf("Expected channel to be %d, but it is %d", want.type_, header.FrameType)
			}

			payload, err := parseFrameBody(&buf)
			if err != nil {
				t.Fatalf("%+v", err)
			}
			if !testEqual(want.body, payload) {
				t.Errorf("Roundtrip produced different results:\n %s", testDiff(want.body, payload))
			}
		})
	}
}

func BenchmarkFrameMarshal(b *testing.B) {
	for _, tt := range exampleFrames {
		b.Run(tt.label, func(b *testing.B) {
			b.ReportAllocs()
			var buf buffer

			for i := 0; i < b.N; i++ {
				err := writeFrame(&buf, tt.frame)
				if err != nil {
					b.Error(fmt.Sprintf("%+v", err))
				}
				bytesSink = buf.bytes()
				buf.reset()
			}
		})
	}
}
func BenchmarkFrameUnmarshal(b *testing.B) {
	for _, tt := range exampleFrames {
		b.Run(tt.label, func(b *testing.B) {
			b.ReportAllocs()
			var buf buffer
			err := writeFrame(&buf, tt.frame)
			if err != nil {
				b.Error(fmt.Sprintf("%+v", err))
			}
			data := buf.bytes()
			buf.reset()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				buf := &buffer{b: data}
				_, err := parseFrameHeader(buf)
				if err != nil {
					b.Errorf("%+v", err)
				}

				_, err = parseFrameBody(buf)
				if err != nil {
					b.Errorf("%+v", err)
				}
			}
		})
	}
}

var bytesSink []byte

func BenchmarkMarshal(b *testing.B) {
	for _, typ := range allTypes {
		b.Run(fmt.Sprintf("%T", typ), func(b *testing.B) {
			b.ReportAllocs()
			var buf buffer

			for i := 0; i < b.N; i++ {
				err := marshal(&buf, typ)
				if err != nil {
					b.Error(fmt.Sprintf("%+v", err))
				}
				bytesSink = buf.bytes()
				buf.reset()
			}
		})
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	for _, type_ := range allTypes {
		b.Run(fmt.Sprintf("%T", type_), func(b *testing.B) {
			var buf buffer
			err := marshal(&buf, type_)
			if err != nil {
				b.Error(fmt.Sprintf("%+v", err))
			}
			data := buf.bytes()
			newType := reflect.New(reflect.TypeOf(type_)).Interface()

			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				err = unmarshal(&buffer{b: data}, newType)
				if err != nil {
					b.Error(fmt.Sprintf("%v", err))
				}
			}
		})
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	_, updateFuzzCorpus := os.LookupEnv("UPDATE_FUZZ_CORPUS")

	for _, type_ := range allTypes {
		t.Run(fmt.Sprintf("%T", type_), func(t *testing.T) {
			var buf buffer
			err := marshal(&buf, type_)
			if err != nil {
				t.Fatal(fmt.Sprintf("%+v", err))
			}

			if updateFuzzCorpus {
				name := fmt.Sprintf("%T.bin", type_)
				name = strings.TrimPrefix(name, "amqp.")
				name = strings.TrimPrefix(name, "*amqp.")
				path := filepath.Join("fuzz/marshal/corpus", name)
				err = ioutil.WriteFile(path, buf.bytes(), 0644)
				if err != nil {
					t.Error(err)
				}
			}

			// handle special case around nil type
			if type_ == nil {
				err = unmarshal(&buf, nil)
				if err != nil {
					t.Fatal(fmt.Sprintf("%+v", err))
					return
				}
				return
			}

			newType := reflect.New(reflect.TypeOf(type_))
			err = unmarshal(&buf, newType.Interface())
			if err != nil {
				t.Fatal(fmt.Sprintf("%+v", err))
				return
			}
			cmpType := reflect.Indirect(newType).Interface()
			if !testEqual(type_, cmpType) {
				t.Errorf("Roundtrip produced different results:\n %s", testDiff(type_, cmpType))
			}
		})
	}
}

func TestReadAny(t *testing.T) {
	for _, type_ := range generalTypes {
		t.Run(fmt.Sprintf("%T", type_), func(t *testing.T) {
			var buf buffer
			err := marshal(&buf, type_)
			if err != nil {
				t.Errorf("%+v", err)
			}

			got, err := readAny(&buf)
			if err != nil {
				t.Fatalf("%+v", err)
			}

			if !testEqual(type_, got) {
				t.Errorf("Roundtrip produced different results:\n %s", testDiff(type_, got))
			}
		})
	}
}

var (
	allTypes = append(protoTypes, generalTypes...)

	protoTypes = []interface{}{
		&performOpen{
			ContainerID:         "foo",
			Hostname:            "bar.host",
			MaxFrameSize:        4200,
			ChannelMax:          13,
			OutgoingLocales:     []symbol{"fooLocale"},
			IncomingLocales:     []symbol{"barLocale"},
			OfferedCapabilities: []symbol{"fooCap"},
			DesiredCapabilities: []symbol{"barCap"},
			Properties: map[symbol]interface{}{
				"fooProp": int32(45),
			},
		},
		&performBegin{
			RemoteChannel:       4321,
			NextOutgoingID:      730000,
			IncomingWindow:      9876654,
			OutgoingWindow:      123555,
			HandleMax:           9757,
			OfferedCapabilities: []symbol{"fooCap"},
			DesiredCapabilities: []symbol{"barCap"},
			Properties: map[symbol]interface{}{
				"fooProp": int32(45),
			},
		},
		&performAttach{
			Name:               "fooName",
			Handle:             435982,
			Role:               roleSender,
			SenderSettleMode:   sndSettle(ModeMixed),
			ReceiverSettleMode: rcvSettle(ModeSecond),
			Source: &source{
				Address:      "fooAddr",
				Durable:      DurabilityUnsettledState,
				ExpiryPolicy: ExpiryLinkDetach,
				Timeout:      635,
				Dynamic:      true,
				DynamicNodeProperties: map[symbol]interface{}{
					"lifetime-policy": deleteOnClose,
				},
				DistributionMode: "some-mode",
				Filter: filter{
					"foo:filter": &describedType{
						descriptor: "foo:filter",
						value:      "bar value",
					},
				},
				Outcomes:     []symbol{"amqp:accepted:list"},
				Capabilities: []symbol{"barCap"},
			},
			Target: &target{
				Address:      "fooAddr",
				Durable:      DurabilityUnsettledState,
				ExpiryPolicy: ExpiryLinkDetach,
				Timeout:      635,
				Dynamic:      true,
				DynamicNodeProperties: map[symbol]interface{}{
					"lifetime-policy": deleteOnClose,
				},
				Capabilities: []symbol{"barCap"},
			},
			Unsettled: unsettled{
				"fooDeliveryTag": &stateAccepted{},
			},
			IncompleteUnsettled:  true,
			InitialDeliveryCount: 3184,
			MaxMessageSize:       75983,
			OfferedCapabilities:  []symbol{"fooCap"},
			DesiredCapabilities:  []symbol{"barCap"},
			Properties: map[symbol]interface{}{
				"fooProp": int32(45),
			},
		},
		role(true),
		&unsettled{
			"fooDeliveryTag": &stateAccepted{},
		},
		&source{
			Address:      "fooAddr",
			Durable:      DurabilityUnsettledState,
			ExpiryPolicy: ExpiryLinkDetach,
			Timeout:      635,
			Dynamic:      true,
			DynamicNodeProperties: map[symbol]interface{}{
				"lifetime-policy": deleteOnClose,
			},
			DistributionMode: "some-mode",
			Filter: filter{
				"foo:filter": &describedType{
					descriptor: "foo:filter",
					value:      "bar value",
				},
			},
			Outcomes:     []symbol{"amqp:accepted:list"},
			Capabilities: []symbol{"barCap"},
		},
		&target{
			Address:      "fooAddr",
			Durable:      DurabilityUnsettledState,
			ExpiryPolicy: ExpiryLinkDetach,
			Timeout:      635,
			Dynamic:      true,
			DynamicNodeProperties: map[symbol]interface{}{
				"lifetime-policy": deleteOnClose,
			},
			Capabilities: []symbol{"barCap"},
		},
		&performFlow{
			NextIncomingID: uint32Ptr(354),
			IncomingWindow: 4352,
			NextOutgoingID: 85324,
			OutgoingWindow: 24378634,
			Handle:         uint32Ptr(341543),
			DeliveryCount:  uint32Ptr(31341),
			LinkCredit:     uint32Ptr(7634),
			Available:      uint32Ptr(878321),
			Drain:          true,
			Echo:           true,
			Properties: map[symbol]interface{}{
				"fooProp": int32(45),
			},
		},
		&performTransfer{
			Handle:             34983,
			DeliveryID:         uint32Ptr(564),
			DeliveryTag:        []byte("foo tag"),
			MessageFormat:      uint32Ptr(34),
			Settled:            true,
			More:               true,
			ReceiverSettleMode: rcvSettle(ModeSecond),
			State:              &stateReceived{},
			Resume:             true,
			Aborted:            true,
			Batchable:          true,
			Payload:            []byte("very important payload"),
		},
		&performDisposition{
			Role:      roleSender,
			First:     5644444,
			Last:      uint32Ptr(423),
			Settled:   true,
			State:     &stateReleased{},
			Batchable: true,
		},
		&performDetach{
			Handle: 4352,
			Closed: true,
			Error: &Error{
				Condition:   ErrorNotAllowed,
				Description: "foo description",
				Info: map[string]interface{}{
					"other": "info",
					"and":   uint16(875),
				},
			},
		},
		&performDetach{
			Handle: 4352,
			Closed: true,
			Error: &Error{
				Condition:   ErrorLinkRedirect,
				Description: "",
				// payload is bigger than map8 encoding size
				Info: map[string]interface{}{
					"hostname":     "redirected.myservicebus.example.org",
					"network-host": "redirected.myservicebus.example.org",
					"port":         uint32(5671),
					"address":      "amqps://redirected.myservicebus.example.org:5671/path",
				},
			},
		},
		ErrorCondition("the condition"),
		&Error{
			Condition:   ErrorNotAllowed,
			Description: "foo description",
			Info: map[string]interface{}{
				"other": "info",
				"and":   uint16(875),
			},
		},
		&performEnd{
			Error: &Error{
				Condition:   ErrorNotAllowed,
				Description: "foo description",
				Info: map[string]interface{}{
					"other": "info",
					"and":   uint16(875),
				},
			},
		},
		&performClose{
			Error: &Error{
				Condition:   ErrorNotAllowed,
				Description: "foo description",
				Info: map[string]interface{}{
					"other": "info",
					"and":   uint16(875),
				},
			},
		},
		&Message{
			Header: &MessageHeader{
				Durable:       true,
				Priority:      234,
				TTL:           10 * time.Second,
				FirstAcquirer: true,
				DeliveryCount: 32,
			},
			DeliveryAnnotations: Annotations{
				int64(42): "answer",
			},
			Annotations: Annotations{
				int64(42): "answer",
			},
			Properties: &MessageProperties{
				MessageID:          "yo",
				UserID:             []byte("baz"),
				To:                 "me",
				Subject:            "sup?",
				ReplyTo:            "you",
				CorrelationID:      uint64(34513),
				ContentType:        "text/plain",
				ContentEncoding:    "UTF-8",
				AbsoluteExpiryTime: time.Date(2018, 01, 13, 14, 24, 07, 0, time.UTC),
				CreationTime:       time.Date(2018, 01, 13, 14, 14, 07, 0, time.UTC),
				GroupID:            "fooGroup",
				GroupSequence:      89324,
				ReplyToGroupID:     "barGroup",
			},
			ApplicationProperties: map[string]interface{}{
				"baz": "foo",
			},
			Data: [][]byte{
				[]byte("A nice little data payload."),
				[]byte("More payload."),
			},
			Value: uint8(42),
			Footer: Annotations{
				"hash": []uint8{0, 1, 2, 34, 5, 6, 7, 8, 9, 0},
			},
		},
		&MessageHeader{
			Durable:       true,
			Priority:      234,
			TTL:           10 * time.Second,
			FirstAcquirer: true,
			DeliveryCount: 32,
		},
		&MessageProperties{
			MessageID:          "yo",
			UserID:             []byte("baz"),
			To:                 "me",
			Subject:            "sup?",
			ReplyTo:            "you",
			CorrelationID:      uint64(34513),
			ContentType:        "text/plain",
			ContentEncoding:    "UTF-8",
			AbsoluteExpiryTime: time.Date(2018, 01, 13, 14, 24, 07, 0, time.UTC),
			CreationTime:       time.Date(2018, 01, 13, 14, 14, 07, 0, time.UTC),
			GroupID:            "fooGroup",
			GroupSequence:      89324,
			ReplyToGroupID:     "barGroup",
		},
		&stateReceived{
			SectionNumber: 234,
			SectionOffset: 8973,
		},
		&stateAccepted{},
		&stateRejected{
			Error: &Error{
				Condition:   ErrorStolen,
				Description: "foo description",
				Info: map[string]interface{}{
					"other": "info",
					"and":   int32(uint16(875)),
				},
			},
		},
		&stateReleased{},
		&stateModified{
			DeliveryFailed:    true,
			UndeliverableHere: true,
			MessageAnnotations: Annotations{
				"more": "annotations",
			},
		},
		lifetimePolicy(typeCodeDeleteOnClose),
		SenderSettleMode(1),
		ReceiverSettleMode(1),
		&saslInit{
			Mechanism:       "FOO",
			InitialResponse: []byte("BAR\x00RESPONSE\x00"),
			Hostname:        "me",
		},
		&saslMechanisms{
			Mechanisms: []symbol{"FOO", "BAR", "BAZ"},
		},
		&saslOutcome{
			Code:           codeSASLSysPerm,
			AdditionalData: []byte("here's some info for you..."),
		},
		milliseconds(10 * time.Second),
		symbol("a symbol"),
		map[symbol]interface{}{
			"hash": []uint8{0, 1, 2, 34, 5, 6, 7, 8, 9, 0},
		},
	}

	generalTypes = []interface{}{
		nil,
		UUID{1, 2, 3, 4, 5, 6, 7, 8, 10, 11, 12, 13, 14, 15, 16},
		bool(true),
		int8(math.MaxInt8),
		int8(math.MinInt8),
		int16(math.MaxInt16),
		int16(math.MinInt16),
		int32(math.MaxInt32),
		int32(math.MinInt32),
		int64(math.MaxInt64),
		int64(math.MinInt64),
		uint8(math.MaxUint8),
		uint16(math.MaxUint16),
		uint32(math.MaxUint32),
		uint64(math.MaxUint64),
		float32(math.Pi),
		float32(-math.Pi),
		float32(math.NaN()),
		float32(-math.NaN()),
		float64(math.Pi),
		float64(-math.Pi),
		float64(math.NaN()),
		float64(-math.NaN()),
		describedType{
			descriptor: binary.BigEndian.Uint64([]byte{0x00, 0x00, 0x46, 0x8C, 0x00, 0x00, 0x00, 0x04}),
			value:      "amqp.annotation.x-opt-offset > '312'",
		},
		map[interface{}]interface{}{
			int32(-1234): []uint8{0, 1, 2, 34, 5, 6, 7, 8, 9, 0},
		},
		map[string]interface{}{
			"hash": []uint8{0, 1, 2, 34, 5, 6, 7, 8, 9, 0},
		},
		ArrayUByte{1, 2, 3, math.MaxUint8, 0},
		[]int8{1, 2, 3, math.MaxInt8, math.MinInt8},
		[]uint16{1, 2, 3, math.MaxUint16, 0},
		[]uint16{1, 2, 3, math.MaxInt8, 0},
		[]int16{1, 2, 3, math.MaxInt16, math.MinInt16},
		[]int16{1, 2, 3, math.MaxInt8, math.MinInt8},
		[]uint32{1, 2, 3, math.MaxUint32, 0},
		[]uint32{1, 2, 3, math.MaxUint8, 0},
		[]int32{1, 2, 3, math.MaxInt32, math.MinInt32},
		[]int32{1, 2, 3, math.MaxInt8, math.MinInt8},
		[]uint64{1, 2, 3, math.MaxUint64, 0},
		[]uint64{1, 2, 3, math.MaxUint8, 0},
		[]int64{1, 2, 3, math.MaxInt64, math.MinInt64},
		[]int64{1, 2, 3, math.MaxInt8, math.MinInt8},
		[]float32{math.Pi, -math.Pi, float32(math.NaN()), float32(-math.NaN())},
		[]float64{math.Pi, -math.Pi, math.NaN(), -math.NaN()},
		[]bool{true, false, true, false},
		[]string{"FOO", "BAR", "BAZ"},
		[]symbol{"FOO", "BAR", "BAZ"},
		[][]byte{[]byte("FOO"), []byte("BAR"), []byte("BAZ")},
		[]time.Time{time.Date(2018, 01, 27, 16, 16, 59, 0, time.UTC)},
		[]UUID{
			{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
			{16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 31},
		},
		[]interface{}{int16(1), "hello", false},
	}
)

func sndSettle(m SenderSettleMode) *SenderSettleMode {
	return &m
}
func rcvSettle(m ReceiverSettleMode) *ReceiverSettleMode {
	return &m
}

func uint32Ptr(u uint32) *uint32 {
	return &u
}
