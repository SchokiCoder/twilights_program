// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

use std::io;
use std::io::Read;

fn run_program()
{
	todo!();
}

fn main()
{
	let mut input: [u8; 1] = [0];
	let mut stdin = io::stdin();

	println!("run program? (y/n)");
	match stdin.read(&mut input) {
	Ok(_) => {}
	
	Err(_) => {
		panic!("Cannot read stdin");
	}}

	match input[0] as char {
	'y' => {
		run_program();
	}
	
	_ => {}
	}
}
