package main

import (
	"log"
	"os"

	"github.com/jozsefsallai/unityweb"
	"github.com/urfave/cli/v2"
)

func unpack(c *cli.Context) error {
	input := c.String("input")
	output := c.String("output")

	pkg, err := unityweb.FromPackageFile(input)
	if err != nil {
		return err
	}

	err = pkg.Dump(output)
	if err != nil {
		return err
	}

	return nil
}

func pack(c *cli.Context) error {
	input := c.String("input")
	output := c.String("output")

	pkg, err := unityweb.PackDirectory(input)
	if err != nil {
		return err
	}

	err = pkg.PackToFile(output)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	app := &cli.App{
		Name:    "unityweb",
		Usage:   "A tool for unpacking and repacking Unity Web data files.",
		Version: "1.0.2",
		Commands: []*cli.Command{
			{
				Name:    "unpack",
				Aliases: []string{"u", "un", "extract", "x"},
				Usage:   "Unpack a Unity Web data file into a directory.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Aliases:  []string{"i"},
						Usage:    "The path to the Unity Web data file.",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "output",
						Aliases:  []string{"o"},
						Usage:    "The path to the output directory.",
						Required: true,
					},
				},
				Action: unpack,
			},
			{
				Name:    "pack",
				Aliases: []string{"p", "repack", "r"},
				Usage:   "Pack a directory into a Unity Web data file.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Aliases:  []string{"i"},
						Usage:    "The path to the input directory.",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "output",
						Aliases:  []string{"o"},
						Usage:    "The path to the output Unity Web data file.",
						Required: true,
					},
				},
				Action: pack,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
