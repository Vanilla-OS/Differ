package types

/*
 * 	License: GPL-3.0-or-later
 * 	Authors:
 * 		Mateus Melchiades <matbme@duck.com>
 * 	Copyright: 2023
 */

import "time"

type Package struct {
	Name, Version string
}

type Release struct {
	Digest   string
	Date     time.Time
	Packages []Package
}
