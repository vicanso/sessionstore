package store

import (
	"math/rand"
	"time"

	"github.com/vicanso/cod"
	"github.com/vicanso/cookies"
	"github.com/vicanso/keygrip"
)

type (
	// Options session store options
	Options struct {
		// session的缓存有效期
		TTL time.Duration
		// session id 的key（cookie的名字或http头）
		Key string
		// cookie options
		CookieOptions *cookies.Options
		// 签名使用的密钥
		SignKeys    []string
		IDGenerator func() string
	}
	// Store session store
	Store struct {
		opts    *Options
		ctx     *cod.Context
		cookies *cookies.Cookies
	}
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// generateID gen id
func generateID() string {
	b := make([]rune, 24)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (s *Store) getCookies() *cookies.Cookies {
	if s.cookies != nil {
		return s.cookies
	}
	opts := s.opts
	cookiesOptions := opts.CookieOptions
	cookiesOptions.Keys = opts.SignKeys
	if s.cookies == nil {
		s.cookies = cookies.New(s.ctx, cookiesOptions)
	}
	return s.cookies
}

func (s *Store) getIDFromCookies() string {
	opts := s.opts
	ck := s.getCookies()
	signed := len(opts.SignKeys) != 0
	return ck.Get(opts.Key, signed)
}

// GetID GetID
func (s *Store) GetID() string {
	opts := s.opts
	if opts.CookieOptions != nil {
		return s.getIDFromCookies()
	}
	ctx := s.ctx
	id := ctx.GetRequestHeader(opts.Key)
	if len(opts.SignKeys) == 0 {
		return id
	}
	kg := keygrip.New(opts.SignKeys)
	signValue := ctx.GetRequestHeader(opts.Key + cookies.SigSuffix)
	if kg.Verify([]byte(id), []byte(signValue)) {
		return id
	}
	return ""
}

func (s *Store) setIDToCookies(id string) error {
	opts := s.opts
	signed := len(opts.SignKeys) != 0
	ck := s.getCookies()
	return ck.Set(ck.CreateCookie(opts.Key, id), signed)
}

// CreateID create id
func (s *Store) CreateID() (id string, err error) {
	opts := s.opts
	if opts.IDGenerator != nil {
		id = opts.IDGenerator()
	} else {
		id = generateID()
	}
	if opts.CookieOptions != nil {
		err = s.setIDToCookies(id)
		return
	}
	ctx := s.ctx
	ctx.SetHeader(opts.Key, id)
	if len(opts.SignKeys) == 0 {
		return
	}
	kg := keygrip.New(opts.SignKeys)
	signValue := string(kg.Sign([]byte(id)))
	ctx.SetHeader(opts.Key+cookies.SigSuffix, signValue)
	return
}

// SetOptions set options
func (s *Store) SetOptions(opts *Options) {
	s.opts = opts
}

// GetTTL get ttl of session
func (s *Store) GetTTL() time.Duration {
	return s.opts.TTL
}
