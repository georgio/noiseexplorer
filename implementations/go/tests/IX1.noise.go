/*
IX1:
  -> e, s
  <- e, ee, se, s
  -> es
  <-
  ->
*/

/* ---------------------------------------------------------------- *
 * PARAMETERS                                                       *
 * ---------------------------------------------------------------- */

package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/binary"
	"encoding/hex"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/hkdf"
	"hash"
	"io"
	"math"
)

/* ---------------------------------------------------------------- *
 * TYPES                                                            *
 * ---------------------------------------------------------------- */

type keypair struct {
	pk [32]byte
	sk [32]byte
}

type messagebuffer struct {
	ne         [32]byte
	ns         []byte
	ciphertext []byte
}

type cipherstate struct {
	k [32]byte
	n uint64
}

type symmetricstate struct {
	cs cipherstate
	ck [32]byte
	h  [32]byte
}

type handshakestate struct {
	ss  symmetricstate
	s   keypair
	e   keypair
	rs  [32]byte
	re  [32]byte
	psk [32]byte
	i   bool
}

type noisesession struct {
	hs  handshakestate
	h   [32]byte
	cs1 cipherstate
	cs2 cipherstate
	mc  uint64
}

/* ---------------------------------------------------------------- *
 * CONSTANTS                                                        *
 * ---------------------------------------------------------------- */

var emptyKey = [32]byte{
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00,
}

var minNonce = uint64(0)
/* ---------------------------------------------------------------- *
 * UTILITY FUNCTIONS                                                *
 * ---------------------------------------------------------------- */

func getPublicKey(kp keypair) [32]byte {
	return kp.pk
}

func isEmptyKey(k [32]byte) bool {
	return subtle.ConstantTimeCompare(k[:], emptyKey[:]) == 1
}

/* ---------------------------------------------------------------- *
 * PRIMITIVES                                                       *
 * ---------------------------------------------------------------- */

func incrementNonce(n uint64) uint64 {
	return n + 1
}

func dh(sk [32]byte, pk [32]byte) [32]byte {
	var ss [32]byte
	curve25519.ScalarMult(&ss, &sk, &pk)
	return ss
}

func generateKeypair() keypair {
	var pk [32]byte
	var sk [32]byte
	_, _ = rand.Read(sk[:])
	curve25519.ScalarBaseMult(&pk, &sk)
	return keypair{pk, sk}
}

func generatePublicKey(sk [32]byte) [32]byte {
	var pk [32]byte
	curve25519.ScalarBaseMult(&pk, &sk)
	return pk
}

func encrypt(k [32]byte, n uint64, ad []byte, plaintext []byte) []byte {
	var nonce [12]byte
	var ciphertext []byte
	enc, _ := chacha20poly1305.New(k[:])
	binary.LittleEndian.PutUint64(nonce[4:], n)
	ciphertext = enc.Seal(nil, nonce[:], plaintext, ad)
	return ciphertext
}

func decrypt(k [32]byte, n uint64, ad []byte, ciphertext []byte) (bool, []byte, []byte) {
	var nonce [12]byte
	var plaintext []byte
	enc, err := chacha20poly1305.New(k[:])
	binary.LittleEndian.PutUint64(nonce[4:], n)
	plaintext, err = enc.Open(nil, nonce[:], ciphertext, ad)
	return (err == nil), ad, plaintext
}

func getHash(a []byte, b []byte) [32]byte {
	return blake2s.Sum256(append(a, b...))
}

func hashProtocolName(protocolName []byte) [32]byte {
	var h [32]byte
	if len(protocolName) <= 32 {
		for i := 0; i < 32; i++ {
			if i < len(protocolName) {
				h[i] = protocolName[i]
			} else {
				h[i] = byte(0x00)
			}
		}
	} else {
		h = getHash(protocolName, []byte{})
	}
	return h
}

func blake2HkdfInterface() hash.Hash {
	h, _ := blake2s.New256([]byte{})
	return h
}

func getHkdf(ck [32]byte, ikm []byte) ([32]byte, [32]byte, [32]byte) {
	var k1 [32]byte
	var k2 [32]byte
	var k3 [32]byte
	output := hkdf.New(blake2HkdfInterface, ikm[:], ck[:], []byte{})
	io.ReadFull(output, k1[:])
	io.ReadFull(output, k2[:])
	io.ReadFull(output, k3[:])
	return k1, k2, k3
}
/* ---------------------------------------------------------------- *
 * STATE MANAGEMENT                                                 *
 * ---------------------------------------------------------------- */

/* CipherState */
func initializeKey(k [32]byte) cipherstate {
	return cipherstate{k, minNonce}
}

func hasKey(cs cipherstate) bool {
	return !isEmptyKey(cs.k)
}

func setNonce(cs cipherstate, newNonce uint64) cipherstate {
	return cipherstate{cs.k, newNonce}
}

func encryptWithAd(cs cipherstate, ad []byte, plaintext []byte) (cipherstate, []byte) {
	e := encrypt(cs.k, cs.n, ad, plaintext)
	csi := setNonce(cs, incrementNonce(cs.n))
	return csi, e
}

func decryptWithAd(cs cipherstate, ad []byte, ciphertext []byte) (cipherstate, []byte, bool) {
	valid, ad, plaintext := decrypt(cs.k, cs.n, ad, ciphertext)
	csi := setNonce(cs, incrementNonce(cs.n))
	return csi, plaintext, valid
}

func reKey(cs cipherstate) cipherstate {
	var ki [32]byte
	e := encrypt(cs.k, math.MaxUint64, []byte{}, emptyKey[:])
	for i := 0; i < 32; i++ {
		ki[i] = e[i]
	}
	return cipherstate{ki, cs.n}
}

/* SymmetricState */

func initializeSymmetric(protocolName []byte) symmetricstate {
	h := hashProtocolName(protocolName)
	ck := h
	cs := initializeKey(emptyKey)
	return symmetricstate{cs, ck, h}
}

func mixKey(ss symmetricstate, ikm [32]byte) symmetricstate {
	ck, tempK, _ := getHkdf(ss.ck, ikm[:])
	csi := initializeKey(tempK)
	return symmetricstate{csi, ck, ss.h}
}

func mixHash(ss symmetricstate, data []byte) symmetricstate {
	return symmetricstate{ss.cs, ss.ck, getHash(ss.h[:], data)}
}

func mixKeyAndHash(ss symmetricstate, ikm [32]byte) symmetricstate {
	ck, tempH, tempK := getHkdf(ss.ck, ikm[:])
	ssi := mixHash(symmetricstate{ss.cs, ck, ss.h}, tempH[:])
	return symmetricstate{initializeKey(tempK), ck, ssi.h}
}

func getHandshakeHash(ss symmetricstate) [32]byte {
	return ss.h
}

func encryptAndHash(ss symmetricstate, plaintext []byte) (symmetricstate, []byte) {
	var csi cipherstate
	var ciphertext []byte
	if hasKey(ss.cs) {
		csi, ciphertext = encryptWithAd(ss.cs, ss.h[:], plaintext)
	} else {
		csi, ciphertext = ss.cs, plaintext
	}
	ssi := mixHash(symmetricstate{csi, ss.ck, ss.h}, ciphertext)
	return ssi, ciphertext
}

func decryptAndHash(ss symmetricstate, ciphertext []byte) (symmetricstate, []byte, bool) {
	var csi cipherstate
	var plaintext []byte
	var valid bool
	if hasKey(ss.cs) {
		csi, plaintext, valid = decryptWithAd(ss.cs, ss.h[:], ciphertext)
	} else {
		csi, plaintext, valid = ss.cs, ciphertext, true
	}
	ssi := mixHash(symmetricstate{csi, ss.ck, ss.h}, ciphertext)
	return ssi, plaintext, valid
}

func split(ss symmetricstate) (cipherstate, cipherstate) {
	tempK1, tempK2, _ := getHkdf(ss.ck, []byte{})
	cs1 := initializeKey(tempK1)
	cs2 := initializeKey(tempK2)
	return cs1, cs2
}

/* HandshakeState */

func initializeInitiator(prologue []byte, s keypair, rs [32]byte, psk [32]byte) handshakestate {
	var ss symmetricstate
	var e keypair
	var re [32]byte
	name := []byte("Noise_IX1_25519_ChaChaPoly_BLAKE2s")
	ss = mixHash(initializeSymmetric(name), prologue)
	return handshakestate{ss, s, e, rs, re, psk, true}
}

func initializeResponder(prologue []byte, s keypair, rs [32]byte, psk [32]byte) handshakestate {
	var ss symmetricstate
	var e keypair
	var re [32]byte
	name := []byte("Noise_IX1_25519_ChaChaPoly_BLAKE2s")
	ss = mixHash(initializeSymmetric(name), prologue)
	return handshakestate{ss, s, e, rs, re, psk, false}
}

func writeMessageA(hs handshakestate, payload []byte) (handshakestate, messagebuffer) {
	ss, s, e, rs, re, psk, initiator := hs.ss, hs.s, hs.e, hs.rs, hs.re, hs.psk, hs.i
	ne, ns, ciphertext := emptyKey, []byte{}, []byte{}
	esk, _ := hex.DecodeString("893e28b9dc6ca8d611ab664754b8ceb7bac5117349a4439a6b0569da977c464a")
	copy(e.sk[:], esk[:])
	e.pk = generatePublicKey(e.sk)
	ne = e.pk
	ss = mixHash(ss, ne[:])
	/* No PSK, so skipping mixKey */
	ss, ns = encryptAndHash(ss, s.pk[:])
	ss, ciphertext = encryptAndHash(ss, payload)
	hs = handshakestate{ss, s, e, rs, re, psk, initiator}
	messageBuffer := messagebuffer{ne, ns, ciphertext}
	return hs, messageBuffer
}

func writeMessageB(hs handshakestate, payload []byte) (handshakestate, messagebuffer) {
	ss, s, e, rs, re, psk, initiator := hs.ss, hs.s, hs.e, hs.rs, hs.re, hs.psk, hs.i
	ne, ns, ciphertext := emptyKey, []byte{}, []byte{}
	esk, _ := hex.DecodeString("bbdb4cdbd309f1a1f2e1456967fe288cadd6f712d65dc7b7793d5e63da6b375b")
	copy(e.sk[:], esk[:])
	e.pk = generatePublicKey(e.sk)
	ne = e.pk
	ss = mixHash(ss, ne[:])
	/* No PSK, so skipping mixKey */
	ss = mixKey(ss, dh(e.sk, re))
	ss = mixKey(ss, dh(e.sk, rs))
	ss, ns = encryptAndHash(ss, s.pk[:])
	ss, ciphertext = encryptAndHash(ss, payload)
	hs = handshakestate{ss, s, e, rs, re, psk, initiator}
	messageBuffer := messagebuffer{ne, ns, ciphertext}
	return hs, messageBuffer
}

func writeMessageC(hs handshakestate, payload []byte) ([32]byte, messagebuffer, cipherstate, cipherstate) {
	ss, s, e, rs, re, psk, initiator := hs.ss, hs.s, hs.e, hs.rs, hs.re, hs.psk, hs.i
	ne, ns, ciphertext := emptyKey, []byte{}, []byte{}
	ss = mixKey(ss, dh(e.sk, rs))
	ss, ciphertext = encryptAndHash(ss, payload)
	hs = handshakestate{ss, s, e, rs, re, psk, initiator}
	messageBuffer := messagebuffer{ne, ns, ciphertext}
	cs1, cs2 := split(ss)
	return hs.ss.h, messageBuffer, cs1, cs2
}

func writeMessageRegular(cs cipherstate, payload []byte) (cipherstate, messagebuffer) {
	/* No handshakestate */
	ne, ns, ciphertext := emptyKey, []byte{}, []byte{}
	cs, ciphertext = encryptWithAd(cs, []byte{}, payload)
	messageBuffer := messagebuffer{ne, ns, ciphertext}
	return cs, messageBuffer
}

func readMessageA(hs handshakestate, message messagebuffer) (handshakestate, []byte, bool) {
	ss, s, e, rs, re, psk, initiator := hs.ss, hs.s, hs.e, hs.rs, hs.re, hs.psk, hs.i
	valid1 := true
	re = message.ne
	ss = mixHash(ss, re[:])
	/* No PSK, so skipping mixKey */
	ss, ns, valid1 := decryptAndHash(ss, message.ns)
	if !valid1 || len(ns) != 32 {
		return hs, []byte{}, false
	}
	for i := 0; i < 32; i++ { rs[i] = ns[i] }
	ss, plaintext, valid2 := decryptAndHash(ss, message.ciphertext)
	if !valid2 {
		return hs, []byte{}, false
	}
	hs = handshakestate{ss, s, e, rs, re, psk, initiator}
	return hs, plaintext, (valid1 && valid2)
}

func readMessageB(hs handshakestate, message messagebuffer) (handshakestate, []byte, bool) {
	ss, s, e, rs, re, psk, initiator := hs.ss, hs.s, hs.e, hs.rs, hs.re, hs.psk, hs.i
	valid1 := true
	re = message.ne
	ss = mixHash(ss, re[:])
	/* No PSK, so skipping mixKey */
	ss = mixKey(ss, dh(e.sk, re))
	ss = mixKey(ss, dh(s.sk, re))
	ss, ns, valid1 := decryptAndHash(ss, message.ns)
	if !valid1 || len(ns) != 32 {
		return hs, []byte{}, false
	}
	for i := 0; i < 32; i++ { rs[i] = ns[i] }
	ss, plaintext, valid2 := decryptAndHash(ss, message.ciphertext)
	if !valid2 {
		return hs, []byte{}, false
	}
	hs = handshakestate{ss, s, e, rs, re, psk, initiator}
	return hs, plaintext, (valid1 && valid2)
}

func readMessageC(hs handshakestate, message messagebuffer) ([32]byte, []byte, bool, cipherstate, cipherstate) {
	ss, s, e, rs, re, psk, initiator := hs.ss, hs.s, hs.e, hs.rs, hs.re, hs.psk, hs.i
	valid1 := true
	ss = mixKey(ss, dh(s.sk, re))
	ss, plaintext, valid2 := decryptAndHash(ss, message.ciphertext)
	if !valid2 {
		return emptyKey, []byte{}, false, hs.ss.cs, hs.ss.cs
	}
	hs = handshakestate{ss, s, e, rs, re, psk, initiator}
	cs1, cs2 := split(ss)
	return hs.ss.h, plaintext, (valid1 && valid2), cs1, cs2
}

func readMessageRegular(cs cipherstate, message messagebuffer) (cipherstate, []byte, bool) {
	/* No handshakestate */
	/* No encrypted keys */
	csi, plaintext, valid2 := decryptWithAd(cs, []byte{}, message.ciphertext)
	if !valid2 {
		return cs, []byte{}, false
	}
	return csi, plaintext, valid2
}



/* ---------------------------------------------------------------- *
 * PROCESSES                                                        *
 * ---------------------------------------------------------------- */

func InitSession(initiator bool, prologue []byte, s keypair, rs [32]byte) noisesession {
	var session noisesession
	psk := emptyKey
	if initiator {
		session.hs = initializeInitiator(prologue, s, rs, psk)
	} else {
		session.hs = initializeResponder(prologue, s, rs, psk)
	}
	session.mc = 0
	return session
}

func SendMessage(session noisesession, message []byte) (noisesession, messagebuffer) {
	var hs handshakestate
	var messageBuffer messagebuffer
	hs = session.hs
	if session.mc == 0 {
		hs, messageBuffer = writeMessageA(hs, message)
	}
	if session.mc == 1 {
		hs, messageBuffer = writeMessageB(hs, message)
	}
	if session.mc == 2 {
		session.h, messageBuffer, session.cs1, session.cs2 = writeMessageC(hs, message)
		session.hs = handshakestate{}
	}
	if session.mc > 2 {
		if hs.i {
			session.cs1, messageBuffer = writeMessageRegular(session.cs1, message)
		} else {
			session.cs2, messageBuffer = writeMessageRegular(session.cs2, message)
		}
	}
	session.mc = session.mc + 1
	session.hs = hs
	return session, messageBuffer
}

func RecvMessage(session noisesession, message messagebuffer) (noisesession, []byte, bool) {
	var hs handshakestate
	var plaintext []byte
	var valid bool
	hs = session.hs
	if session.mc == 0 {
		hs, plaintext, valid = readMessageA(hs, message)
	}
	if session.mc == 1 {
		hs, plaintext, valid = readMessageB(hs, message)
	}
	if session.mc == 2 {
		session.h, plaintext, valid, session.cs1, session.cs2 = readMessageC(hs, message)
		session.hs = handshakestate{}
	}
	if session.mc > 2 {
		if hs.i {
			session.cs2, plaintext, valid = readMessageRegular(session.cs2, message)
		} else {
			session.cs1, plaintext, valid = readMessageRegular(session.cs1, message)
		}
	}
	session.mc = session.mc + 1
	session.hs = hs
	return session, plaintext, valid
}

func main() {
	prologue, _ := hex.DecodeString("4a6f686e2047616c74")
	var initStatic keypair
	initStaticSk, _ := hex.DecodeString("e61ef9919cde45dd5f82166404bd08e38bceb5dfdfded0a34c8df7ed542214d1")
	copy(initStatic.sk[:], initStaticSk[:])
	initStatic.pk = generatePublicKey(initStatic.sk)
	var respStatic keypair
	respStaticSk, _ := hex.DecodeString("4a3acbfdb163dec651dfa3194dece676d437029c62a408b4c5ea9114246e4893")
	copy(respStatic.sk[:], respStaticSk[:])
	respStatic.pk = generatePublicKey(respStatic.sk)
	initiatorSession := InitSession(true, prologue, initStatic, emptyKey)
	responderSession := InitSession(false, prologue, respStatic, emptyKey)
	payloadA, _ := hex.DecodeString("4c756477696720766f6e204d69736573")
	payloadB, _ := hex.DecodeString("4d757272617920526f746862617264")
	payloadC, _ := hex.DecodeString("462e20412e20486179656b")
	payloadD, _ := hex.DecodeString("4361726c204d656e676572")
	payloadE, _ := hex.DecodeString("4a65616e2d426170746973746520536179")
	payloadF, _ := hex.DecodeString("457567656e2042f6686d20766f6e2042617765726b")
	initiatorSession, messageA := SendMessage(initiatorSession, payloadA)
	responderSession, _, validA := RecvMessage(responderSession, messageA)
	responderSession, messageB := SendMessage(responderSession, payloadB)
	initiatorSession, _, validB := RecvMessage(initiatorSession, messageB)
	initiatorSession, messageC := SendMessage(initiatorSession, payloadC)
	responderSession, _, validC := RecvMessage(responderSession, messageC)
	responderSession, messageD := SendMessage(responderSession, payloadD)
	initiatorSession, _, validD := RecvMessage(initiatorSession, messageD)
	initiatorSession, messageE := SendMessage(initiatorSession, payloadE)
	responderSession, _, validE := RecvMessage(responderSession, messageE)
	responderSession, messageF := SendMessage(responderSession, payloadF)
	initiatorSession, _, validF := RecvMessage(initiatorSession, messageF)
	if validA && validB && validC && validD && validE && validF {
		println("Sanity check PASS for IX1_25519_ChaChaPoly_BLAKE2s.")
	} else {
		println("Sanity check FAIL for IX1_25519_ChaChaPoly_BLAKE2s.")
	}
	cA := hex.EncodeToString(messageA.ne[:]) + hex.EncodeToString(messageA.ns) + hex.EncodeToString(messageA.ciphertext)
	cB := hex.EncodeToString(messageB.ne[:]) + hex.EncodeToString(messageB.ns) + hex.EncodeToString(messageB.ciphertext)
	cC := hex.EncodeToString(messageC.ns) + hex.EncodeToString(messageC.ciphertext)
	cD := hex.EncodeToString(messageD.ns) + hex.EncodeToString(messageD.ciphertext)
	cE := hex.EncodeToString(messageE.ns) + hex.EncodeToString(messageE.ciphertext)
	cF := hex.EncodeToString(messageF.ns) + hex.EncodeToString(messageF.ciphertext)
	tA := "ca35def5ae56cec33dc2036731ab14896bc4c75dbb07a61f879f8e3afa4c79446bc3822a2aa7f4e6981d6538692b3cdf3e6df9eea6ed269eb41d93c22757b75a4c756477696720766f6e204d69736573"
	tB := "95ebc60d2b1fa672c1f46a8aa265ef51bfe38e7ccb39ec5be34069f1448088430363f728c82edd54d4f71350d2c4f7f5af1bed4b49ef675e794191fbe71850f61cecca65b10d1d23042946d173110d1829d651eb05dd70c2669887e79a3a63ab7d965042b982c015e735e4daa3a3e6"
	tC := "cab545be46f12f865c4dfa1bedfad43b3fd74bbba68b782e00a604"
	tD := "8faa9cae6553c0ba4f954103cd73b85dd1c656d180a15225147bbd"
	tE := "b0ac0e3dfef53d59c6ac663edc510022e5a2bc02472b5b572f90a8384838037c1d"
	tF := "90d7cd3c583f94f09f776abb4fa99b430e4f442c710194314a42949d0a167df36a7714ad1b"
	if tA == cA {
		println("Test 1: PASS")
	} else {
		println("Test 1: FAIL")
		println("Expected:	", tA) 
		println("Actual:		", cA)
	}
	if tB == cB {
		println("Test 2: PASS")
	} else {
		println("Test 2: FAIL")
		println("Expected:	", tB) 
		println("Actual:		", cB)
	}
	if tC == cC {
		println("Test 3: PASS")
	} else {
		println("Test 3: FAIL")
		println("Expected:	", tC) 
		println("Actual:		", cC)
	}
	if tD == cD {
		println("Test 4: PASS")
	} else {
		println("Test 4: FAIL")
		println("Expected:	", tD) 
		println("Actual:		", cD)
	}
	if tE == cE {
		println("Test 5: PASS")
	} else {
		println("Test 5: FAIL")
		println("Expected:	", tE) 
		println("Actual:		", cE)
	}
	if tF == cF {
		println("Test 6: PASS")
	} else {
		println("Test 6: FAIL")
		println("Expected:	", tF) 
		println("Actual:		", cF)
	}
}