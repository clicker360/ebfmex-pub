package sess

import (
    "appengine"
    "appengine/datastore"
    "net/http"
	"crypto/md5"
    "time"
	"fmt"
    "io"
)

type Sess struct {
	Md5			string
	Id			string
	User		string
	Name		string
	Expiration	time.Time
	ErrMsg		string `datastore:"-"`
	ErrClass	string `datastore:"-"`
}

func SetSess(w http.ResponseWriter, c appengine.Context, key *datastore.Key, email string, name string) (string, *datastore.Key, error) {
	h := md5.New()
	io.WriteString(h, key.Encode())
	io.WriteString(h, fmt.Sprintf("%s", time.Now()))
	md5 := fmt.Sprintf("%x", h.Sum(nil))
	ex := time.Now().AddDate(0,0,5)
	s := Sess{
		Md5:		md5,
		Id:			key.Encode(),
		User:		email,
		Name:		name,
		Expiration:	ex,
	}
	cKey, err := datastore.Put(c, datastore.NewKey(c, "Sess", email, 0, nil), &s)
	if err != nil {
		return "", nil, err
	}
	// Se crean 2 cookies, una con el key de sesión otra con el número random llave 
	//csc := http.Cookie{ Name: "ebfmex-pub-sesscontrol-ua", Value: md5, Expires: ex, Path: "/" }
	//http.SetCookie(w, &csc)
	//csc = http.Cookie{ Name: "ebfmex-pub-sessid-ua", Value: cKey.Encode(), Expires: ex, Path: "/" }
	//http.SetCookie(w, &csc)
	w.Header().Add("Set-Cookie", fmt.Sprintf("ebfmex-pub-sesscontrol-ua=%s; expires=%s; path=/;", md5, ex.Format("Wed, 07-Oct-2012 14:23:42 GMT")))
	w.Header().Add("Set-Cookie", fmt.Sprintf("ebfmex-pub-sessid-ua=%s; expires=%s; path=/;", cKey.Encode(), ex.Format("Wed, 07-Oct-2012 14:23:42 GMT")))
	return md5, cKey, err
}

func IsSess(w http.ResponseWriter, r *http.Request, c appengine.Context) (Sess, bool) {
	var s Sess
	if ck, err := r.Cookie("ebfmex-pub-sessid-ua"); err == nil {
		key, _ := datastore.DecodeKey(ck.Value)
		if err := datastore.Get(c, key, &s); err != nil {
			// no hay sesión
			//http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			/* se verifica el control de sesion */
			if cr, err := r.Cookie("ebfmex-pub-sesscontrol-ua"); err == nil {
				if s.Md5 != cr.Value {
					// Md5 no coincide, intenta entrar con otra cookie
					return s, false
				} else if time.Now().After(s.Expiration) {
					// Sesión expirada
					return s, false
				}
				// OK
				// Hay sesión
				return s, true
			}
		}
	}
	return s, false
}


