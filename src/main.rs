// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

use iced::Color;
use iced::{Element, Sandbox, Settings};
use iced::font::*;
use iced::time;
use iced::theme::{self, Theme};
use iced::widget::{column, container, text};
use std::{io};
use std::io::Read;

struct Twiggie {
	theme: Theme,
	start: time::Instant,
	uptime: time::Duration,
}

impl Sandbox for Twiggie {
	type Message = ();

	fn new() -> Twiggie
	{
		let palette = theme::Palette {
			background: Color::from_rgb(0.24, 0.643, 0.565),
			text: Color::from_rgb(0.212, 0.996, 0.8),
			primary: Color::from_rgb(0.5, 0.5, 0.0),
			success: Color::from_rgb(0.0, 1.0, 0.0),
			danger: Color::from_rgb(1.0, 0.0, 0.0),
		};

		let ret = Twiggie {
			theme: Theme::custom(palette),
			start: time::Instant::now(),
			uptime: time::Duration::from_nanos(0),
		};

		return ret;
	}

	fn view(&self) -> Element<Self::Message>
	{
		let intro = match self.uptime.as_secs() { 
		0 | 1 => {
			text("YOU ARE NOW")
		}
		
		_ => {
			text("").size(0)
		}};

		let content = column![
			self.theme,
			intro
		];

		container(content)
	}
	
	fn update(&mut self, _msg: Self::Message)
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
	let mut settings = Settings::default();
	settings.antialiasing = false;
	settings.default_font.stretch = Stretch::Normal;
	settings.default_font.monospaced = true;
	settings.default_text_size = 14.0;
	settings.window.size = (320, 200);
	let settings = settings;

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
		match Twiggie::run(settings) {
		Ok(_) => {}
		
		Err(e) => {
			panic!("Could not run program: {}", e);
		}}
	}
	
	_ => {}
	}
}
