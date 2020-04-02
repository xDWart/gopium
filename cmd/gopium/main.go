package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"1pkg/gopium"
	"1pkg/gopium/runners"

	"github.com/spf13/cobra"
)

// list of all cli vars
var (
	// cli command iteself
	cli *cobra.Command
	// target platform vars
	tcompiler, tarch string
	tcpulines        []int
	// package parser vars
	pname, ppath    string
	pbenvs, pbflags []string
	// walker strategies vars
	wname, wregex   string
	wdeep, wbackref bool
	tagtype         string
	// global vars
	timeout int
	// global running context
	gctx    context.Context
	gcancel func()
)

// init cli command runner
// and global running context
func init() {
	// set root cli command app
	cli = &cobra.Command{
		Use:     "gopium -w walker_name -n package_name strategy_name#1 strategy_name#2 strategy_name#3 ...",
		Short:   gopium.STAMP,
		Version: gopium.VERSION,
		Example: "gopium -w json_std -n 1pkg/gopium -g soft filter_pads memory_pack separate_padding_cpu_l1_top separate_padding_cpu_l1_bottom",
		Long: `
Gopium is the tool which was designed to automate and simplify non trivial actions on structs, like:
 - cpu cache alignment
 - memory packing
 - false sharing guarding
 - auto annotation
 - generic fields management
 - other relevant activities

In order to use gopium cli you need to provide at least victim package_name list of strategies which will be applied one by one and walker_name.
Outcome of execution is fully defined by list of strategies and walker_name combination. List of strategies modifies victim structs in package.
Walker facilitates and insures that outcome is written to one of supported destination (check walker_name flag).

Gopium supports next strategies list: 
 - process_tag_group (uses gopium fields tags annotation in order to process different set of strategies on different groups and then combine results in single struct result)

 - memory_pack (rearranges structure fields to obtain optimal memory utilization)
 - memory_unpack (rearranges structure field list to obtain inflated memory utilization)
	
 - cache_rounding_cpu_l1 (fits structure into cpu cache line #1 by adding bottom rounding cpu cache padding)
 - cache_rounding_cpu_l2 (fits structure into cpu cache line #2 by adding bottom rounding cpu cache padding)
 - cache_rounding_cpu_l3 (fits structure into cpu cache line #3 by adding bottom rounding cpu cache padding)

 - false_sharing_cpu_l1 (guards structure from false sharing by adding extra cpu cache line #1 paddings for each structure field)
 - false_sharing_cpu_l2 (guards structure from false sharing by adding extra cpu cache line #1 paddings for each structure field)
 - false_sharing_cpu_l3 (guards structure from false sharing by adding extra cpu cache line #1 paddings for each structure field)

 - separate_padding_system_alignment_top (separates structure with extra system alignment padding by adding the padding at the top)
 - separate_padding_cpu_l1_top (separates structure with extra cpu cache line #1 padding by adding the padding at the top)
 - separate_padding_cpu_l2_top (separates structure with extra cpu cache line #2 padding by adding the padding at the top)
 - separate_padding_cpu_l3_top (separates structure with extra cpu cache line #3 padding by adding the padding at the top)
 - separate_padding_system_alignment_bottom (separates structure with extra system alignment padding by adding the padding at the bottom)
 - separate_padding_cpu_l1_bottom (separates structure with extra cpu cache line #1 padding by adding the padding at the bottom)
 - separate_padding_cpu_l2_bottom (separates structure with extra cpu cache line #2 padding by adding the padding at the bottom)
 - separate_padding_cpu_l3_bottom (separates structure with extra cpu cache line #3 padding by adding the padding at the bottom)

 - explicit_padings_system_alignment (explicitly aligns each structure field to system alignment padding by adding missing paddings for each field)
 - explicit_padings_type_natural (explicitly aligns each structure field to max type alignment padding by adding missing paddings for each field)

 - doc_fields_annotate (adds size doc annotation for each structure field and aggregated size annotation for whole structure)
 - comment_fields_annotate (adds size comment annotation for each structure field and aggregated size annotation for whole structure)
 - doc_struct_stamp (adds doc gopium stamp to structure)
 - comment_struct_stamp (adds comment gopium stamp to structure)

 - name_lexicographical_ascending (sorts fields accordingly to their names in ascending order)
 - name_lexicographical_descending (sorts fields accordingly to their names descending order)
 - name_length_ascending (sorts fields accordingly to their names length in ascending order)
 - name_length_descending (sorts fields accordingly to their names length in descending order)
 - type_lexicographical_ascending (sorts fields accordingly to their types in ascending order)
 - type_lexicographical_descending (sorts fields accordingly to their types in descending order)
 - type_length_ascending (sorts fields accordingly to their types length in ascending order)
 - type_length_descending (sorts fields accordingly to their types length in descending order)

 - embedded_ascending (sorts fields accordingly to their embeded flag in ascending order)
 - embedded_descending (sorts fields accordingly to their embeded flag in descending order)
 - exported_ascending (sorts fields accordingly to their exported flag in ascending order)
 - exported_descending (sorts fields accordingly to their exported flag in descending order)

 - filter_pads (filters out all structure padding fields)
 - filter_embedded (filters out all structure embedded fields)
 - filter_not_embedded (filters out all structure not embedded fields)
 - filter_exported (filters out all structure exported fields)
 - filter_not_exported (filters out all structure not exported fields)
 - remove_tag_group (removes gopium fields tags annotation)

 - nope (does nothing by returning original structure)
 - void (does nothing by returning void struct)

Notes:
 - it might be useful to use filter_pads in pipes with other strategies to clean paddings first
 - process_tag_group currently supports only next fields tags annotation formats:
  - gopium:"stg,stg,stg" processed as default group
  - gopium:"group:def;stg,stg,stg" processed as named group
- by specifying tag_type you can automatically generate fields tags annotation suitable for process_tag_group
		`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, stgs []string) error {
			return runners.NewCliApp(
				tcompiler,
				tarch,
				tcpulines,
				pname,
				ppath,
				pbenvs,
				pbflags,
				wname,
				wregex,
				wdeep,
				wbackref,
				stgs,
				tagtype,
				timeout,
			).Run(cmd.Context())
		},
	}
	// set target_compiler flag
	cli.Flags().StringVarP(
		&tcompiler,
		"target_compiler",
		"c",
		"gc",
		`
Target platform compiler name, possible values are:
 - gc
 - gccgo
		`,
	)
	// set target_architecture flag
	cli.Flags().StringVarP(
		&tarch,
		"target_architecture",
		"a",
		"amd64",
		`
Target platform architecture name, possible values are: 
 - 386
 - arm
 - arm64
 - amd64
 - mips
 - etc.
		`,
	)
	// set target_cpu_cache_line_sizes flag
	cli.Flags().IntSliceVarP(
		&tcpulines,
		"target_cpu_cache_line_sizes",
		"l",
		[]int{64, 64, 64},
		`
Target platform CPU cache line sizes in bytes, cache line size is set one by one l1,l2,l3,...
For now only 3 lines of cache are supported by strategies.
		`,
	)
	// set required package_name flag
	cli.Flags().StringVarP(
		&pname,
		"package_name",
		"n",
		"",
		"Go package name, full package name is expected.",
	)
	cli.MarkFlagRequired("package_name")
	// set package_path flag
	cli.Flags().StringVarP(
		&ppath,
		"package_path",
		"p",
		"",
		"Go package path, path to root of the package is expected.",
	)
	// set package_build_envs flag
	cli.Flags().StringSliceVarP(
		&pbenvs,
		"package_build_envs",
		"",
		[]string{},
		"Go package build envs, additional list of building envs is expected.",
	)
	// set package_build_flags flag
	cli.Flags().StringSliceVarP(
		&pbflags,
		"package_build_flags",
		"",
		[]string{},
		"Go package build flags, additional list of building flags is expected.",
	)
	// set required walker_name flag
	cli.Flags().StringVarP(
		&wname,
		"walker_name",
		"w",
		"",
		`
Gopium walker name, possible values are: 
 - json_std (prints json encoded result to stdout)
 - xml_std (prints xml encoded result to stdout)
 - csv_std (prints csv encoded result to stdout)
 - json_files (prints json encoded result to files inside package directory)
 - xml_files (prints xml encoded result to files inside package directory)
 - csv_files (prints csv encoded result to files inside package directory)
 - sync_ast (directly syncs result as go code in orinal package)
		`,
	)
	cli.MarkFlagRequired("walker_name")
	// set walker_regexp flag
	cli.Flags().StringVarP(
		&wregex,
		"walker_regexp",
		"r",
		".*",
		`
Gopium walker regexp, regexp that defines which structures would be visited.
Visiting is done only if structure name matches the regexp.
		`,
	)
	// set walker_deep flag
	cli.Flags().BoolVarP(
		&wdeep,
		"walker_deep",
		"d",
		true,
		`
Gopium walker deep flag, flag that defines type of nested scopes visiting.
By default it visits all nested scopes.
		`,
	)
	// set walker_backref flag
	cli.Flags().BoolVarP(
		&wbackref,
		"walker_backref",
		"b",
		true,
		`
Gopium walker backref flag, flag that defines type of names referencing.
By default any previous visited types have affect on future relevant visits.
		`,
	)
	// set tag_type flag
	cli.Flags().StringVarP(
		&tagtype,
		"tag_type",
		"g",
		"none",
		`
Gopium strategy tag type write policy, possible values are: 
 - none (do nothing with tag)
 - soft (write tag only if not exists yet)
 - force (overwrite tag even if tag exists already)
		`,
	)
	// set timeout flag
	cli.Flags().IntVarP(
		&timeout,
		"timeout",
		"t",
		0,
		"Global timeout of cli command in seconds, considered only if value > 0.",
	)
	// prepare global context
	// with cancelation
	// on system signals
	gctx, gcancel = context.WithCancel(context.Background())
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, os.Kill)
		select {
		case <-gctx.Done():
		case <-sig:
			gcancel()
		}
	}()
}

// main gopium cli entry point
func main() {
	// execute cobra cli command
	// on global running context
	// and log error if any
	defer gcancel()
	if err := cli.ExecuteContext(gctx); err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
}
