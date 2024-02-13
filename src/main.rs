// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

use iced::{Element, Sandbox, Settings};
use iced::time;
use iced::widget::{text};
use std::{io};
use std::io::Read;

/*#[derive(Debug, Clone, Copy)]
pub enum Message {
	Click
}*/

struct Twiggie {
	//message: Message,
	start: time::Instant,
	uptime: time::Duration,
}

impl Sandbox for Twiggie {
	type Message = ();

	fn new() -> Twiggie
	{
		return Twiggie {
			start: time::Instant::now(),
			uptime: time::Duration::from_nanos(0),
		};
	}

	fn view(&self) -> Element<Self::Message>
	{
		let intro = match self.uptime.as_secs() { 
		0 | 1 => {
			text("YOU ARE NOW").size(50)
		}
		
		_ => {
			text("").size(0)
		}};

		intro.into()
	}
	
	fn update(&mut self, msg: Self::Message)
	{	
		self.uptime = self.start.elapsed();
	}
	
	fn title(&self) -> String
	{
		"Twilights program".into()
	}
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
		Twiggie::run(Settings::default());
	}
	
	_ => {}
	}
}
