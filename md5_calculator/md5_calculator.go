package md5_calculator

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Calculator struct {
	Client *http.Client
}

func (c Calculator) CalcMD5Hash(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := c.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %s", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read resp body: %s", err)
	}

	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:]), nil
}
