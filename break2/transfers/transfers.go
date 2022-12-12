package transfers

import (
	"encoding/gob"
	"log"
	"net"
	"zero/break/cryptos"
)

const BUFSIZE = 32 * 1024

type Transfers struct {
	cp cryptos.AesCrypt
}

func NewTransfers(key string) (Transfers, error) {
	c := cryptos.NewAesCrypt()
	err := c.SetKey(key)
	if err != nil {
		return Transfers{}, err
	}
	return Transfers{
		cp: c,
	}, nil
}

func (c *Transfers) EnData(s []byte) ([]byte, error) {
	b, err := c.cp.EnCode(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *Transfers) DeData(s []byte) ([]byte, error) {
	b, err := c.cp.DeCode(s)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *Transfers) DataCopyEn(dst, src net.Conn) {
	buf := make([]byte, BUFSIZE)
	enc := gob.NewEncoder(dst)
	for {
		// read data
		nr, er := src.Read(buf)
		if nr > 0 {
			// after read encrypt data
			b, err := c.cp.EnCode(buf[:nr])
			if err != nil {
				log.Printf("<Transfer> encrypt data error: %v", err)
				break
			}
			// write data
			ew := enc.Encode(b)
			if ew != nil {
				log.Printf("<Transfer> write crypto data error: %v", ew)
				break
			}
		}
		if er != nil {
			log.Printf("<Transfer> read data error: %v", er)
			break
		}
	}
}

func (c *Transfers) DataCopyDe(dst, src net.Conn) {
	buf := make([]byte, BUFSIZE)
	dec := gob.NewDecoder(src)
	for {
		// read data
		err := dec.Decode(&buf)
		if err != nil {
			log.Printf("<Transfer> read crypto data error: %v", err)
			break
		}

		// handle date
		b, err := c.cp.DeCode(buf)
		if err != nil {
			log.Printf("<Transfer> decrypt data error: %v", err)
			break
		}

		// write data
		nw, err := dst.Write(b)
		if err != nil || nw != len(b) {
			log.Printf("<Transfer> write data error: %v", err)
			break
		}
	}
}
