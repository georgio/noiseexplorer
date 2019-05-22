/*
NNpsk2:
  -> e
  <- e, ee, psk
  ->
  <-
*/

/* ---------------------------------------------------------------- *
 * PARAMETERS                                                       *
 * ---------------------------------------------------------------- */

#[macro_use]
pub(crate) mod macros;

pub(crate) mod prims;
pub(crate) mod state;

pub mod consts;
pub mod error;
pub mod noisesession;
pub mod types;