#![allow(non_snake_case, non_upper_case_globals, unused_assignments, unused_imports)]

use noiseexplorer_nx::{
	consts::{DHLEN, MAC_LENGTH},
	error::NoiseError,
	noisesession::NoiseSession,
	types::{Keypair, PrivateKey, PublicKey},
};

fn decode_str(s: &str) -> Vec<u8> {
 	hex::decode(s).unwrap()
 }

#[test]
fn noiseexplorer_test_nx() {

	let mut prologue: Vec<u8> = Vec::new();
	// length = message length + mac length + ne length (32) 
	let mut messageA: Vec<u8> = Vec::new();
	// length = message length + mac length +
	let mut messageB: Vec<u8> = Vec::new();
	// length = message length + mac length +
	let mut messageC: Vec<u8> = Vec::new();
	// length = message length + mac length +
	let mut messageD: Vec<u8> = Vec::new();
	// length = message length + mac length +
	let mut messageE: Vec<u8> = Vec::new();
	// length = message length + mac length +
	let mut messageF: Vec<u8> = Vec::new();

	prologue = decode_str("4a6f686e2047616c74");
	let initiator_static_private = PrivateKey::from_str("0000000000000000000000000000000000000000000000000000000000000001").unwrap();
	let responder_static_private = PrivateKey::from_str("4a3acbfdb163dec651dfa3194dece676d437029c62a408b4c5ea9114246e4893").unwrap();
	let initiator_static_kp = Keypair::from_private_key(initiator_static_private).unwrap();
	let responder_static_kp = Keypair::from_private_key(responder_static_private).unwrap();
	
	let mut initiator_session: NoiseSession = NoiseSession::init_session(true, &prologue[..], initiator_static_kp);
	let mut responder_session: NoiseSession = NoiseSession::init_session(false, &prologue[..], responder_static_kp);
	let initiator_ephemeral_private = PrivateKey::from_str("893e28b9dc6ca8d611ab664754b8ceb7bac5117349a4439a6b0569da977c464a").unwrap();
	let initiator_ephemeral_kp = Keypair::from_private_key(initiator_ephemeral_private).unwrap();
	initiator_session.set_ephemeral_keypair(initiator_ephemeral_kp);
	let responder_ephemeral_private = PrivateKey::from_str("bbdb4cdbd309f1a1f2e1456967fe288cadd6f712d65dc7b7793d5e63da6b375b").unwrap();
	let responder_ephemeral_kp = Keypair::from_private_key(responder_ephemeral_private).unwrap();
	responder_session.set_ephemeral_keypair(responder_ephemeral_kp);
	messageA.extend_from_slice(&[0u8; DHLEN][..]);
	messageA.extend_from_slice(&decode_str("4c756477696720766f6e204d69736573")[..]);
	let tA: Vec<u8> = Vec::from(&decode_str("ca35def5ae56cec33dc2036731ab14896bc4c75dbb07a61f879f8e3afa4c79444c756477696720766f6e204d69736573")[..]);
	// messageA length is 32 + payload length,
	// payload starts at index 32
	initiator_session.send_message(&mut messageA[..]).unwrap();
	responder_session.recv_message(&mut messageA.clone()[..]).unwrap();
	messageB.extend_from_slice(&[0u8; DHLEN][..]);
	messageB.extend_from_slice(&[0u8; DHLEN+MAC_LENGTH][..]);
	messageB.extend_from_slice(&decode_str("4d757272617920526f746862617264")[..]);
	messageB.extend_from_slice(&[0u8; MAC_LENGTH][..]);
	let tB: Vec<u8> = Vec::from(&decode_str("95ebc60d2b1fa672c1f46a8aa265ef51bfe38e7ccb39ec5be34069f1448088431b7ab475ba0987fba04b749be49e6b43fe538cfca25a1c591a7ed09f19c9b9e7d042761a2fd2762cf2cb2062ce2c61253452b8383eb2ddc9ba2237b96d97b4e866ba73f55165a736ad03e68594ce25")[..]);
	// messageB length is 96 + payload length,
	// payload starts at index 80
	responder_session.send_message(&mut messageB[..]).unwrap();
	initiator_session.recv_message(&mut messageB.clone()[..]).unwrap();
	messageC.extend_from_slice(&decode_str("462e20412e20486179656b")[..]);
	messageC.extend_from_slice(&[0u8; MAC_LENGTH][..]);
	let tC: Vec<u8> = Vec::from(&decode_str("5ab8adddb31ab4f1086c55c3f3ed053f4d78eca7aaf7ba09d486f8")[..]);
	// messageC length is 16 + payload length,
	// payload starts at index 0
	initiator_session.send_message(&mut messageC[..]).unwrap();
	responder_session.recv_message(&mut messageC.clone()[..]).unwrap();
	messageD.extend_from_slice(&decode_str("4361726c204d656e676572")[..]);
	messageD.extend_from_slice(&[0u8; MAC_LENGTH][..]);
	let tD: Vec<u8> = Vec::from(&decode_str("f3bbada5c0a4cd615bed55ee18046ad55efc4f30d318c57b4941e1")[..]);
	// messageD length is 16 + payload length,
	// payload starts at index 0
	responder_session.send_message(&mut messageD[..]).unwrap();
	initiator_session.recv_message(&mut messageD.clone()[..]).unwrap();
	messageE.extend_from_slice(&decode_str("4a65616e2d426170746973746520536179")[..]);
	messageE.extend_from_slice(&[0u8; MAC_LENGTH][..]);
	let tE: Vec<u8> = Vec::from(&decode_str("c1372cf03d2727f6b74f656b587735109ebb6159434a40a65e2e6095c12db5f01c")[..]);
	// messageE length is 16 + payload length,
	// payload starts at index 0
	initiator_session.send_message(&mut messageE[..]).unwrap();
	responder_session.recv_message(&mut messageE.clone()[..]).unwrap();
	messageF.extend_from_slice(&decode_str("457567656e2042f6686d20766f6e2042617765726b")[..]);
	messageF.extend_from_slice(&[0u8; MAC_LENGTH][..]);
	let tF: Vec<u8> = Vec::from(&decode_str("de040777d38c7bf60c4b8c0ca730a9526ff067db990848ac33e9e9970b01efdf00bab518d0")[..]);
	// messageF length is 16 + payload length,
	// payload starts at index 0
	responder_session.send_message(&mut messageF[..]).unwrap();
	initiator_session.recv_message(&mut messageF.clone()[..]).unwrap();
	assert!(tA == messageA, "\n\n\nTest A: FAIL\n\nExpected:\n{:X?}\n\nActual:\n{:X?}", tA, messageA);
	assert!(tB == messageB, "\n\n\nTest B: FAIL\n\nExpected:\n{:X?}\n\nActual:\n{:X?}", tB, messageB);
	assert!(tC == messageC, "\n\n\nTest C: FAIL\n\nExpected:\n{:X?}\n\nActual:\n{:X?}", tC, messageC);
	assert!(tD == messageD, "\n\n\nTest D: FAIL\n\nExpected:\n{:X?}\n\nActual:\n{:X?}", tD, messageD);
	assert!(tE == messageE, "\n\n\nTest E: FAIL\n\nExpected:\n{:X?}\n\nActual:\n{:X?}", tE, messageE);
	assert!(tF == messageF, "\n\n\nTest F: FAIL\n\nExpected:\n{:X?}\n\nActual:\n{:X?}", tF, messageF);
}
