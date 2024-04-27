package cmd

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/spf13/cobra"
)

func formatSize(byteLength int) string {
	if byteLength >= 1000*1000*1000 {
		return fmt.Sprintf("%vGB", float64(byteLength)/(1000*1000*1000))
	} else if byteLength >= 1000*1000 {
		return fmt.Sprintf("%vMB", float64(byteLength)/(1000*1000))
	} else if byteLength >= 1000 {
		return fmt.Sprintf("%vKB", float64(byteLength)/1000)
	}
	return fmt.Sprintf("%vB", byteLength)
}

var brotliCompression bool
var gzipCompression bool
var compressionLevel string

var rootCmd = &cobra.Command{
	Use:   "textpress",
	Short: "compresses text files with brotli or gzip",
	Long:  "compresses text files with brotli or gzip first argument, is the input file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("error: no input file was provided")
			os.Exit(1)
		}
		file, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Println("error: invalid input file was provided")
			fmt.Println(err)
		}
		bb := bytes.NewBuffer([]byte{})

		if brotliCompression {
			level := brotli.DefaultCompression
			if compressionLevel == "size" {
				level = brotli.BestCompression
			}
			if compressionLevel == "speed" {
				level = brotli.BestSpeed
			}
			w := brotli.NewWriterV2(bb, level)
			start := time.Now()
			w.Write(file)
			w.Close()
			duration := time.Since(start)
			fmt.Println("Brotli compression:")
			fmt.Println(formatSize(bb.Len()))
			fmt.Printf("In %v\n", duration)
		}
		gb := bytes.NewBuffer([]byte{})

		if gzipCompression {
			level := gzip.DefaultCompression
			if compressionLevel == "size" {
				level = gzip.BestCompression
			}
			if compressionLevel == "speed" {
				level = gzip.BestSpeed
			}
			start := time.Now()
			w, err := gzip.NewWriterLevel(gb, level)
			if err != nil {
				fmt.Println("error: incorrect level")
				os.Exit(1)
			}
			w.Write(file)
			w.Close()
			duration := time.Since(start)
			fmt.Println("GZIP compression:")
			fmt.Println(formatSize(gb.Len()))
			fmt.Printf("In %v\n", duration)
		}

		if len(args) > 1 {
			if gzipCompression && brotliCompression {
				fmt.Println("You selected more than one compression algorithm")
				fmt.Println("1: brotli")
				fmt.Println("2: gzip")
				fmt.Print("Select which you would like to write to the output file: ")
				var opt int
				for _, err := fmt.Scanln(&opt); err != nil; _, err = fmt.Scanln(&opt) {
					fmt.Println("Please enter a valid integer!")
				}
				if opt == 1 {
					os.WriteFile(args[1], bb.Bytes(), 0777)
				} else {
					os.WriteFile(args[1], gb.Bytes(), 0777)

				}

			}
			if brotliCompression {
				os.WriteFile(args[1], bb.Bytes(), 0777)
			} else if gzipCompression {
				os.WriteFile(args[1], gb.Bytes(), 0777)

			}
		}

	},
}

func init() {
	rootCmd.Flags().StringVarP(&compressionLevel, "level", "l", "size", "select compression to best speed or size")
	rootCmd.Flags().BoolVarP(&brotliCompression, "brotli", "b", false, "select brotli compression")
	rootCmd.Flags().BoolVarP(&gzipCompression, "gzip", "g", true, "select gzip compression")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
