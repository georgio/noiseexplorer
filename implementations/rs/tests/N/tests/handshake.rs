#![allow(non_snake_case, non_upper_case_globals)]

use N;

fn decode_str(s: &str) -> Vec<u8> {
    if let Ok(x) = hex::decode(s) {
        x
    } else {
        panic!("{:X?}", hex::decode(s).err());
    }
}

#[test]
fn test() {
    	let prologue = decode_str("4a6f686e2047616c74");
	let initStatic: N::Keypair = N::Keypair::new_k(N::decode_str_32("N::EMPTY_KEY"));
	let respStatic: N::Keypair = N::Keypair::new_k(N::decode_str_32("4a3acbfdb163dec651dfa3194dece676d437029c62a408b4c5ea9114246e4893"));
	let mut initiatorSession: N::NoiseSession =
	N::NoiseSession::InitSession(true, &prologue, initStatic, respStatic.pk.0);
	let mut responderSession: N::NoiseSession =
	N::NoiseSession::InitSession(false, &prologue, respStatic, N::EMPTY_KEY);
	let payloadA = decode_str("4c756477696720766f6e204d69736573");
	let mut messageA: N::MessageBuffer = initiatorSession.SendMessage(&payloadA);
	let mut validA: bool = false;
	if let Some(_x) = responderSession.RecvMessage(&mut messageA) {
	validA = true;
}
	let tA: Vec<u8> = decode_str("ca35def5ae56cec33dc2036731ab14896bc4c75dbb07a61f879f8e3afa4c79441b168ed8bbe8220b52bbbde6593d109d78c299b567f6e69276efcf2659c39073");
	let payloadB = decode_str("4d757272617920526f746862617264");
	let mut messageB: N::MessageBuffer = responderSession.SendMessage(&payloadB);
	let mut validB: bool = false;
	if let Some(_x) = initiatorSession.RecvMessage(&mut messageB) {
	validB = true;
}
	let tB: Vec<u8> = decode_str("a7b5d1962001e9c4d965ea5f133941e9e6989094bcde637a582c34b954f34a");
	let payloadC = decode_str("462e20412e20486179656b");
	let mut messageC: N::MessageBuffer = initiatorSession.SendMessage(&payloadC);
	let mut validC: bool = false;
	if let Some(_x) = responderSession.RecvMessage(&mut messageC) {
	validC = true;
}
	let tC: Vec<u8> = decode_str("16ff2557d5d671abe58c88d2a31b58e3a494ab3a6498124be0ea3f");
	let payloadD = decode_str("4361726c204d656e676572");
	let mut messageD: N::MessageBuffer = responderSession.SendMessage(&payloadD);
	let mut validD: bool = false;
	if let Some(_x) = initiatorSession.RecvMessage(&mut messageD) {
	validD = true;
}
	let tD: Vec<u8> = decode_str("1a6e85b0ef71c38db2c2bf3ebef1d41dc93e26bea6899187d5633d");
	let payloadE = decode_str("4a65616e2d426170746973746520536179");
	let mut messageE: N::MessageBuffer = initiatorSession.SendMessage(&payloadE);
	let mut validE: bool = false;
	if let Some(_x) = responderSession.RecvMessage(&mut messageE) {
	validE = true;
}
	let tE: Vec<u8> = decode_str("00ad2b7d0a03a748d0aefd3accee7bbbcc0bb0ed64d685b2ee8af78997a0245e3f");
	let payloadF = decode_str("457567656e2042f6686d20766f6e2042617765726b");
	let mut messageF: N::MessageBuffer = responderSession.SendMessage(&payloadF);
	let mut validF: bool = false;
	if let Some(_x) = initiatorSession.RecvMessage(&mut messageF) {
	validF = true;
}
	let tF: Vec<u8> = decode_str("5631105c749b9550b27d7926dec0c5b83d4bf207688deccd51b50dd7fc9d5e337bba9c3177");
	if validA && validB && validC && validD && validE && validF {
		println!("Sanity check PASS for N_25519_ChaChaPoly_BLAKE2s.");
	} else {
		println!("Sanity check FAIL for N_25519_ChaChaPoly_BLAKE2s.");
	}
	let mut cA: Vec<u8> = Vec::from(&messageA.ne[..]);
	cA.append(&mut messageA.ns);
	cA.append(&mut messageA.ciphertext);
	let mut cB: Vec<u8> = Vec::from(&messageB.ne[..]);
	cB.append(&mut messageB.ns);
	cB.append(&mut messageB.ciphertext);
	let mut cC: Vec<u8> = messageC.ciphertext;
	let mut cD: Vec<u8> = messageD.ciphertext;
	let mut cE: Vec<u8> = messageE.ciphertext;
	let mut cF: Vec<u8> = messageF.ciphertext;
	if tA == cA {
		println!("Test A: PASS");
	} else {
		println!("Test A: FAIL");
		println!("Expected:	{:X?}", tA);
		println!("Actual:		{:X?}", cA);
	}
	if tB == cB {
		println!("Test B: PASS");
	} else {
		println!("Test B: FAIL");
		println!("Expected:	{:X?}", tB);
		println!("Actual:		{:X?}", cB);
	}
	if tC == cC {
		println!("Test C: PASS");
	} else {
		println!("Test C: FAIL");
		println!("Expected:	{:X?}", tC);
		println!("Actual:		{:X?}", cC);
	}
	if tD == cD {
		println!("Test D: PASS");
	} else {
		println!("Test D: FAIL");
		println!("Expected:	{:X?}", tD);
		println!("Actual:		{:X?}", cD);
	}
	if tE == cE {
		println!("Test E: PASS");
	} else {
		println!("Test E: FAIL");
		println!("Expected:	{:X?}", tE);
		println!("Actual:		{:X?}", cE);
	}
	if tF == cF {
		println!("Test F: PASS");
	} else {
		println!("Test F: FAIL");
		println!("Expected:	{:X?}", tF);
		println!("Actual:		{:X?}", cF);
	}
	assert_eq!(tA, cA);
	assert_eq!(tB, cB);
	assert_eq!(tC, cC);
	assert_eq!(tD, cD);
	assert_eq!(tE, cE);
	assert_eq!(tF, cF);
}