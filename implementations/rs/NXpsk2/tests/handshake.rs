#![allow(non_snake_case, non_upper_case_globals)]

use noiseexplorer_nxpsk2::{noisesession::NoiseSession,
                           types::{Keypair, Message, PrivateKey, Psk, PublicKey}};

#[test]
fn noiseexplorer_test_nxpsk2() {
	if let Ok(prologue,) = Message::from_str("4a6f686e2047616c74",) {
		if let Ok(init_static_private,) = PrivateKey::from_str(
			"0000000000000000000000000000000000000000000000000000000000000001",
		) {
			if let Ok(resp_static_private,) = PrivateKey::from_str(
				"4a3acbfdb163dec651dfa3194dece676d437029c62a408b4c5ea9114246e4893",
			) {
				if let Ok(init_static_kp,) = Keypair::from_private_key(init_static_private,) {
					if let Ok(resp_static_kp,) = Keypair::from_private_key(resp_static_private,) {
						if let Ok(psk,) = Psk::from_str(
							"54686973206973206d7920417573747269616e20706572737065637469766521",
						) {
							let mut initiator_session: NoiseSession = NoiseSession::init_session(
								true,
								prologue.clone(),
								init_static_kp,
								PublicKey::empty(),
								psk.clone(),
							);
							let mut responder_session: NoiseSession = NoiseSession::init_session(
								false,
								prologue,
								resp_static_kp,
								PublicKey::empty(),
								psk,
							);
							if let Ok(initiator_ephemeral_private,) = PrivateKey::from_str(
								"893e28b9dc6ca8d611ab664754b8ceb7bac5117349a4439a6b0569da977c464a",
							) {
								if let Ok(init_ephemeral_kp,) =
									Keypair::from_private_key(initiator_ephemeral_private,)
								{
									initiator_session.set_ephemeral_keypair(init_ephemeral_kp,);
									if let Ok(responder_ephemeral_private,) =
										PrivateKey::from_str("4a6f686e2047616c74",)
									{
										if let Ok(responder_ephemeral_kp,) =
											Keypair::from_private_key(responder_ephemeral_private,)
										{
											responder_session
												.set_ephemeral_keypair(responder_ephemeral_kp,);
											if let Ok(mA,) = Message::from_str(
												"4c756477696720766f6e204d69736573",
											) {
												if let Ok(tA,) = Message::from_str(
													"ca35def5ae56cec33dc2036731ab14896bc4c75dbb07a61f879f8e3afa4c794497bdcf5dc128b7bc5b8f2b6ac1a46dff9f9469337cfac0098f87b2a577cece84",
												) {
													if let Ok(messageA,) =
														initiator_session.send_message(mA,)
													{
														if let Ok(_x,) = responder_session
															.recv_message(messageA.clone(),)
														{
															if let Ok(mB,) = Message::from_str(
																"4d757272617920526f746862617264",
															) {
																if let Ok(tB,) = Message::from_str(
																	"95ebc60d2b1fa672c1f46a8aa265ef51bfe38e7ccb39ec5be34069f144808843c39479e89953f195c89ee9a53f2e291727e15ab09a61b1ca623ee98d3d2549bda7af1881b0ae7ba4bb6e8f71e119927c6c8510ae728cd8c258c6200b71c86e16f934ba80fe35e708f8a52a5e193346",
																) {
																	if let Ok(messageB,) =
																		responder_session
																			.send_message(mB,)
																	{
																		if let Ok(_x,) =
																			initiator_session
																				.recv_message(
																					messageB
																						.clone(),
																				) {
																			if let Ok(mC,) =
																				Message::from_str(
																					"462e20412e20486179656b",
																				) {
																				if let Ok(tC) = Message::from_str("499085038a82c4bc9895c069b9a71ead87545a9184a395d74378e1") {
	if let Ok(messageC) = initiator_session.send_message(mC) {
	if let Ok(_x) = responder_session.recv_message(messageC.clone()) {
	if let Ok(mD) = Message::from_str("4361726c204d656e676572") {
	if let Ok(tD) = Message::from_str("966e81056dafc90a22e1b23039427325cab7791b92bed9a562808e") {
	if let Ok(messageD) = responder_session.send_message(mD) {
	if let Ok(_x) = initiator_session.recv_message(messageD.clone()) {
	if let Ok(mE) = Message::from_str("4a65616e2d426170746973746520536179") {
	if let Ok(tE) = Message::from_str("f5731b781c54f95e5c75bbd6b9b88113de6097618936495b0ba90d545187a3512c") {
	if let Ok(messageE) = initiator_session.send_message(mE) {
	if let Ok(_x) = responder_session.recv_message(messageE.clone()) {
	if let Ok(mF) = Message::from_str("457567656e2042f6686d20766f6e2042617765726b") {
	if let Ok(tF) = Message::from_str("1ad6f2ce261f2f6773363a6f3efc2105c98d960b910629da596e394b052389c66ae988bd84") {
	if let Ok(messageF) = responder_session.send_message(mF) {
	if let Ok(_x) = initiator_session.recv_message(messageF.clone()) {
	assert!(tA == messageA, "\n\n\nTest A: FAIL\n\nExpected:\n{:X?}\n\nActual:\n{:X?}", tA, messageA);
	assert!(tB == messageB, "\n\n\nTest B: FAIL\n\nExpected:\n{:X?}\n\nActual:\n{:X?}", tB, messageB);
	assert!(tC == messageC, "\n\n\nTest C: FAIL\n\nExpected:\n{:X?}\n\nActual:\n{:X?}", tC, messageC);
	assert!(tD == messageD, "\n\n\nTest D: FAIL\n\nExpected:\n{:X?}\n\nActual:\n{:X?}", tD, messageD);
	assert!(tE == messageE, "\n\n\nTest E: FAIL\n\nExpected:\n{:X?}\n\nActual:\n{:X?}", tE, messageE);
	assert!(tF == messageF, "\n\n\nTest F: FAIL\n\nExpected:\n{:X?}\n\nActual:\n{:X?}", tF, messageF);
	}}}}}
	}}}}}}}}}}
																			}
																		}
																	}
																}
															}
														}
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}
