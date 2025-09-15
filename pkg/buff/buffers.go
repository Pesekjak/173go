package buff

import (
	"encoding/binary"
	"io"
)

type MCReader struct {
	r io.Reader
}

func NewReader(r io.Reader) *MCReader {
	return &MCReader{r: r}
}

func (r *MCReader) ReadByte() (byte, error) {
	var data byte
	err := binary.Read(r.r, binary.BigEndian, &data)
	return data, err
}

func (r *MCReader) ReadShort() (int16, error) {
	var data int16
	err := binary.Read(r.r, binary.BigEndian, &data)
	return data, err
}

func (r *MCReader) ReadInt() (int32, error) {
	var data int32
	err := binary.Read(r.r, binary.BigEndian, &data)
	return data, err
}

func (r *MCReader) ReadLong() (int64, error) {
	var data int64
	err := binary.Read(r.r, binary.BigEndian, &data)
	return data, err
}

func (r *MCReader) ReadFloat() (float32, error) {
	var data float32
	err := binary.Read(r.r, binary.BigEndian, &data)
	return data, err
}

func (r *MCReader) ReadDouble() (float64, error) {
	var data float64
	err := binary.Read(r.r, binary.BigEndian, &data)
	return data, err
}

func (r *MCReader) ReadBool() (bool, error) {
	var val byte
	if err := binary.Read(r.r, binary.BigEndian, &val); err != nil {
		return false, err
	}
	return val != 0, nil
}

func (r *MCReader) ReadString8() (string, error) {
	length, err := r.ReadShort()
	if err != nil {
		return "", err
	}
	if length < 0 { // Basic sanity check
		return "", io.ErrUnexpectedEOF
	}

	strBytes := make([]byte, length)
	if _, err := io.ReadFull(r.r, strBytes); err != nil {
		return "", err
	}

	return string(strBytes), nil
}

func (r *MCReader) ReadString16() (string, error) {
	length, err := r.ReadShort()
	if err != nil {
		return "", err
	}
	if length < 0 { // Basic sanity check
		return "", io.ErrUnexpectedEOF
	}

	runes := make([]rune, length)
	for i := int16(0); i < length; i++ {
		var char uint16
		if err := binary.Read(r.r, binary.BigEndian, &char); err != nil {
			return "", err
		}
		runes[i] = rune(char)
	}

	return string(runes), nil
}

type MCWriter struct {
	w io.Writer
}

func NewWriter(w io.Writer) *MCWriter {
	return &MCWriter{w: w}
}

func (w *MCWriter) WriteByte(data byte) error {
	return binary.Write(w.w, binary.BigEndian, data)
}

func (w *MCWriter) WriteShort(data int16) error {
	return binary.Write(w.w, binary.BigEndian, data)
}

func (w *MCWriter) WriteInt(data int32) error {
	return binary.Write(w.w, binary.BigEndian, data)
}

func (w *MCWriter) WriteLong(data int64) error {
	return binary.Write(w.w, binary.BigEndian, data)
}

func (w *MCWriter) WriteFloat(data float32) error {
	return binary.Write(w.w, binary.BigEndian, data)
}

func (w *MCWriter) WriteDouble(data float64) error {
	return binary.Write(w.w, binary.BigEndian, data)
}

func (w *MCWriter) WriteBool(data bool) error {
	var val byte = 0
	if data {
		val = 1
	}
	return binary.Write(w.w, binary.BigEndian, val)
}

func (w *MCWriter) WriteString8(data string) error {
	if err := w.WriteShort(int16(len(data))); err != nil {
		return err
	}
	_, err := io.WriteString(w.w, data)
	return err
}

func (w *MCWriter) WriteString16(data string) error {
	runes := []rune(data)
	if err := w.WriteShort(int16(len(runes))); err != nil {
		return err
	}

	for _, r := range runes {
		if err := binary.Write(w.w, binary.BigEndian, uint16(r)); err != nil {
			return err
		}
	}

	return nil
}
