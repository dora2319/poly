package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/TimothyStiles/poly"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

/******************************************************************************
Oct, 15, 2020

Testing command line utilities via subroutines can be annoying so
if you're doing it from the commandline be sure to compile first.
From the project's root directory often use:

go build && go install && go test -v ./...

To accurately test your commands you MUST make sure to rebuild and reinstall
before you run your tests. Otherwise your system version will be out of
date and will give you results using an older build.

Happy hacking,
Tim


TODO:

write subtest to check for empty output before merge
******************************************************************************/

func TestConvertPipe(t *testing.T) {

	var writeBuffer bytes.Buffer
	app := application()
	app.Writer = &writeBuffer

	args := os.Args[0:1]                                // Name of the program.
	args = append(args, "c", "-i", "gbk", "-o", "json") // Append a flag

	file, _ := ioutil.ReadFile("../data/puc19.gbk")
	app.Reader = bytes.NewReader(file)

	err := app.Run(args)

	if err != nil {
		t.Fatalf("Run error: %s", err)
	}

	// getting test sequence from non-pipe io to compare against io to stdout
	baseTestSequence := poly.ReadGbk("../data/puc19.gbk")

	pipeOutputTestSequence := poly.ParseJSON(writeBuffer.Bytes())

	if diff := cmp.Diff(baseTestSequence, pipeOutputTestSequence, cmpopts.IgnoreFields(poly.Feature{}, "ParentSequence")); diff != "" {
		t.Errorf(" mismatch from convert pipe input test (-want +got):\n%s", diff)
	}

}

func TestConvertFile(t *testing.T) {

	app := application()

	args := os.Args[0:1]                                                                // Name of the program.
	args = append(args, "c", "-o", "json", "../data/puc19.gbk", "../data/t4_intron.gb") // Append a flag

	err := app.Run(args)

	if err != nil {
		t.Fatalf("Run error: %s", err)
	}

	puc19InputTestSequence := poly.ReadGbk("../data/puc19.gbk")
	puc19OutputTestSequence := poly.ReadJSON("../data/puc19.json")

	//clearing test data.
	os.Remove("../data/puc19.json")

	// compared input gff from resulting output json. Fail test and print diff if error.
	if diff := cmp.Diff(puc19InputTestSequence, puc19OutputTestSequence, cmpopts.IgnoreFields(poly.Feature{}, "ParentSequence")); diff != "" {
		t.Errorf(" mismatch from concurrent gbk input test (-want +got):\n%s", diff)
	}

	t4InputTestSequence := poly.ReadGbk("../data/t4_intron.gb")
	t4OutputTestSequence := poly.ReadJSON("../data/t4_intron.json")

	// clearing test data.
	os.Remove("../data/t4_intron.json")

	// compared input gbk from resulting output json. Fail test and print diff if error.
	if diff := cmp.Diff(t4InputTestSequence, t4OutputTestSequence, cmpopts.IgnoreFields(poly.Feature{}, "ParentSequence")); diff != "" {
		t.Errorf(" mismatch from concurrent gbk input test (-want +got):\n%s", diff)
	}
}
func TestHashFile(t *testing.T) {

	puc19GbkBlake3Hash := "4b0616d1b3fc632e42d78521deb38b44fba95cca9fde159e01cd567fa996ceb9"
	var writeBuffer bytes.Buffer

	app := application()
	app.Writer = &writeBuffer

	// testing file matching hash
	args := os.Args[0:1]                             // Name of the program.
	args = append(args, "hash", "../data/puc19.gbk") // Append a flag

	err := app.Run(args)

	if err != nil {
		t.Fatalf("Run error: %s", err)
	}

	if writeBuffer.Len() == 0 {
		t.Error("TestHash did not write output to desired writer.")
	}

	hashOutputString := strings.TrimSpace(writeBuffer.String())
	if hashOutputString != puc19GbkBlake3Hash {
		t.Errorf("TestHashFile has failed. Returned %q, want %q", hashOutputString, puc19GbkBlake3Hash)
	}

}

func TestHashPipe(t *testing.T) {

	puc19GbkBlake3Hash := "4b0616d1b3fc632e42d78521deb38b44fba95cca9fde159e01cd567fa996ceb9"
	var writeBuffer bytes.Buffer

	// create a mock application
	app := application()
	app.Writer = &writeBuffer
	file, _ := ioutil.ReadFile("../data/puc19.gbk")
	app.Reader = bytes.NewReader(file)

	args := os.Args[0:1]                     // Name of the program.
	args = append(args, "hash", "-i", "gbk") // Append a flag

	err := app.Run(args)

	if err != nil {
		t.Fatalf("Run error: %s", err)
	}

	hashOutputString := strings.TrimSpace(writeBuffer.String())

	if hashOutputString != puc19GbkBlake3Hash {
		t.Errorf("TestHashPipe has failed. Returned %q, want %q", hashOutputString, puc19GbkBlake3Hash)
	}

}

func TestHashJSON(t *testing.T) {
	// testing json write output

	puc19GbkBlake3Hash := "4b0616d1b3fc632e42d78521deb38b44fba95cca9fde159e01cd567fa996ceb9"

	app := application()

	args := os.Args[0:1]                                           // Name of the program.
	args = append(args, "hash", "-o", "json", "../data/puc19.gbk") // Append a flag
	err := app.Run(args)

	if err != nil {
		t.Fatalf("Run error: %s", err)
	}

	hashOutputString := poly.ReadJSON("../data/puc19.json").SequenceHash
	os.Remove("../data/puc19.json")

	if hashOutputString != puc19GbkBlake3Hash {
		t.Errorf("TestHashJSON has failed. Returned %q, want %q", hashOutputString, puc19GbkBlake3Hash)
	}

}

func TestHashFunctions(t *testing.T) {

	hashfunctionsTestMap := map[string]string{
		"":            "4b0616d1b3fc632e42d78521deb38b44fba95cca9fde159e01cd567fa996ceb9", // default / defaults to blake3
		"BLAKE2b_256": "51f33a2aaa3ad3884e18822b02ffd37ee27888c7d6103c97bf9628843e1fc445",
		"BLAKE2b_384": "c56ff217ca0e09b1d20a910ec51cb17fd828ae47cf6208bfca2381021cdff86ffd6d0975f6955593400232b1cafd96fa",
		"BLAKE2b_512": "9bc9066ea9f7bb58cc432175a9a99ebb2332247fc97b4f73318882e83bcff647d88503c6bfb02f6d6468af9b88dc303f3270b50922585e669e0c990fb384f4f6",
		"BLAKE2s_256": "3ca14269a42d9ea4c619ced109f879f1bba76469be519b9644ed977753e05353",
		"BLAKE3":      "4b0616d1b3fc632e42d78521deb38b44fba95cca9fde159e01cd567fa996ceb9",
		"MD5":         "2c6d6bdc39c1d59c9aa3d332ddc28190",
		"NO":          "aaaaaaaccaccgctaccagcggtggtttgtttgccggatcaagagctaccaactctttttccgaaggtaactggcttcagcagagcgcagataccaaatactgttcttctagtgtagccgtagttaggccaccacttcaagaactctgtagcaccgcctacatacctcgctctgctaatcctgttaccagtggctgctgccagtggcgataagtcgtgtcttaccgggttggactcaagacgatagttaccggataaggcgcagcggtcgggctgaacggggggttcgtgcacacagcccagcttggagcgaacgacctacaccgaactgagatacctacagcgtgagctatgagaaagcgccacgcttcccgaagggagaaaggcggacaggtatccggtaagcggcagggtcggaacaggagagcgcacgagggagcttccagggggaaacgcctggtatctttatagtcctgtcgggtttcgccacctctgacttgagcgtcgatttttgtgatgctcgtcaggggggcggagcctatggaaaaacgccagcaacgcggcctttttacggttcctggccttttgctggccttttgctcacatgttctttcctgcgttatcccctgattctgtggataaccgtattaccgcctttgagtgagctgataccgctcgccgcagccgaacgaccgagcgcagcgagtcagtgagcgaggaagcggaagagcgcccaatacgcaaaccgcctctccccgcgcgttggccgattcattaatgcagctggcacgacaggtttcccgactggaaagcgggcagtgagcgcaacgcaattaatgtgagttagctcactcattaggcaccccaggctttacactttatgcttccggctcgtatgttgtgtggaattgtgagcggataacaatttcacacaggaaacagctatgaccatgattacgccaagcttgcatgcctgcaggtcgactctagaggatccccgggtaccgagctcgaattcactggccgtcgttttacaacgtcgtgactgggaaaaccctggcgttacccaacttaatcgccttgcagcacatccccctttcgccagctggcgtaatagcgaagaggcccgcaccgatcgcccttcccaacagttgcgcagcctgaatggcgaatggcgcctgatgcggtattttctccttacgcatctgtgcggtatttcacaccgcatatggtgcactctcagtacaatctgctctgatgccgcatagttaagccagccccgacacccgccaacacccgctgacgcgccctgacgggcttgtctgctcccggcatccgcttacagacaagctgtgaccgtctccgggagctgcatgtgtcagaggttttcaccgtcatcaccgaaacgcgcgagacgaaagggcctcgtgatacgcctatttttataggttaatgtcatgataataatggtttcttagacgtcaggtggcacttttcggggaaatgtgcgcggaacccctatttgtttatttttctaaatacattcaaatatgtatccgctcatgagacaataaccctgataaatgcttcaataatattgaaaaaggaagagtatgagtattcaacatttccgtgtcgcccttattcccttttttgcggcattttgccttcctgtttttgctcacccagaaacgctggtgaaagtaaaagatgctgaagatcagttgggtgcacgagtgggttacatcgaactggatctcaacagcggtaagatccttgagagttttcgccccgaagaacgttttccaatgatgagcacttttaaagttctgctatgtggcgcggtattatcccgtattgacgccgggcaagagcaactcggtcgccgcatacactattctcagaatgacttggttgagtactcaccagtcacagaaaagcatcttacggatggcatgacagtaagagaattatgcagtgctgccataaccatgagtgataacactgcggccaacttacttctgacaacgatcggaggaccgaaggagctaaccgcttttttgcacaacatgggggatcatgtaactcgccttgatcgttgggaaccggagctgaatgaagccataccaaacgacgagcgtgacaccacgatgcctgtagcaatggcaacaacgttgcgcaaactattaactggcgaactacttactctagcttcccggcaacaattaatagactggatggaggcggataaagttgcaggaccacttctgcgctcggcccttccggctggctggtttattgctgataaatctggagccggtgagcgtgggtctcgcggtatcattgcagcactggggccagatggtaagccctcccgtatcgtagttatctacacgacggggagtcaggcaactatggatgaacgaaatagacagatcgctgagataggtgcctcactgattaagcattggtaactgtcagaccaagtttactcatatatactttagattgatttaaaacttcatttttaatttaaaaggatctaggtgaagatcctttttgataatctcatgaccaaaatcccttaacgtgagttttcgttccactgagcgtcagaccccgtagaaaagatcaaaggatcttcttgagatcctttttttctgcgcgtaatctgctgcttgcaaac",
		"RIPEMD160":   "43255e2d11394e6429f588d583cd1aa6001cb5d9",
		"SHA1":        "9fe3597d67e31e003bc0aa04054eb37e85e26ffd",
		"SHA244":      "9b674798798e8cfe1b341fe77f575357bb5e4038721ae9fbe41bdd3b",
		"SHA256":      "4a9da025ff81f65a25a73137e5e528ee8923dc47b29a6d12a64989d67aec5f2a",
		"SHA384":      "06f633efe87daf8f903bd1421670f8b4367394ee496408c31a45ae4866a82ec77237ad8f972717e5c0df8bc8bc2db181",
		"SHA3_224":    "2624ef4b8a88e9ad0504a33b6bb236ae9409e83a34e07c76fefbf20b",
		"SHA3_256":    "7ff3f5c0e95705e82701c62ab68e1a6af4591f26b413990eed37d661f4b3ffcb",
		"SHA3_384":    "3f631cbf030e8cbd5c72ec4ba2e775bae4f61c2f5c30c873178f62871d4da413266f908c4aa130ece5412b0243e15101",
		"SHA3_512":    "2f831d14e072044715c4ebfba80920aeefbf86728cf1132afe2688e7fe51a8d7e14203e95da88fd0bcefe064a588cda2497e5be69f304d35088708a0e5f7e843",
		"SHA512":      "66ee5400a01c238d6be1a04e1133f5f2a9fd178ff9c57c5e695b47191aadaddbff5eb885170b3e1715db9c06b1fd4ee88cd599ba98857978b2648e76c7c80ad9",
		"SHA512_224":  "1e543b46774dae76f8a8e3014a214b70b85cb49e7b3dba806b36263a",
		"SHA512_256":  "7691fa37e37ee5a56d0c844d8fde0682713ab4732eeedb973d1a485f29348ac6",
	}

	// sorting keys for easier to read debugging output.
	keys := make([]string, 0, len(hashfunctionsTestMap))
	for k := range hashfunctionsTestMap {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, key := range keys {

		var writeBuffer bytes.Buffer

		app := application()
		app.Writer = &writeBuffer
		args := os.Args[0:1]
		args = append(args, "hash", "-f", key, "../data/puc19.gbk")
		err := app.Run(args)

		if err != nil {
			t.Fatalf("Run error: %s", err)
		}

		if err != nil {
			t.Fatalf("Run error: %s", err)
		}

		if writeBuffer.Len() == 0 {
			t.Error("TestHash did not write output to desired writer.")
		}

		hashOutputString := strings.TrimSpace(writeBuffer.String())
		// fmt.Println("\""+key+"\"", ": ", "\""+hashOutputString+"\",") // <- prints out every
		if hashOutputString != hashfunctionsTestMap[key] {
			t.Errorf("TestHashFunctions for function %q has failed. Returned %q, want %q", key, hashOutputString, hashfunctionsTestMap[key])
		}
	}
}

func TestOptimizeString(t *testing.T) {

	gfpTranslation := "MASKGEELFTGVVPILVELDGDVNGHKFSVSGEGEGDATYGKLTLKFICTTGKLPVPWPTLVTTFSYGVQCFSRYPDHMKRHDFFKSAMPEGYVQERTISFKDDGNYKTRAEVKFEGDTLVNRIELKGIDFKEDGNILGHKLEYNYNSHNVYITADKQKNGIKANFKIRHNIEDGSVQLADHYQQNTPIGDGPVLLPDNHYLSTQSALSKDPNEKRDHMVLLEFVTAAGITHGMDELYK*"
	var writeBuffer bytes.Buffer

	app := application()
	app.Writer = &writeBuffer
	app.Reader = bytes.NewBufferString(gfpTranslation)

	args := os.Args[0:1]                                      // Name of the program.
	args = append(args, "op", "-aa", "-wt", "data/puc19.gbk") // Append a flag
	err := app.Run(args)

	if err != nil {
		t.Fatalf("Run error: %s", err)
	}
	app.Reader = os.Stdin

	// should return codon optimized sequence
	optimizeOutputString := strings.TrimSpace(writeBuffer.String())

	translation := poly.Translate(optimizeOutputString, poly.GetCodonTable(1))

	if translation != gfpTranslation {
		t.Errorf("TestOptimizeCommand for string output has failed. Returned %q, want %q", translation, gfpTranslation)
	}

}

func TestTranslationString(t *testing.T) {
	gfpDnaSequence := "ATGGCTAGCAAAGGAGAAGAACTTTTCACTGGAGTTGTCCCAATTCTTGTTGAATTAGATGGTGATGTTAATGGGCACAAATTTTCTGTCAGTGGAGAGGGTGAAGGTGATGCTACATACGGAAAGCTTACCCTTAAATTTATTTGCACTACTGGAAAACTACCTGTTCCATGGCCAACACTTGTCACTACTTTCTCTTATGGTGTTCAATGCTTTTCCCGTTATCCGGATCATATGAAACGGCATGACTTTTTCAAGAGTGCCATGCCCGAAGGTTATGTACAGGAACGCACTATATCTTTCAAAGATGACGGGAACTACAAGACGCGTGCTGAAGTCAAGTTTGAAGGTGATACCCTTGTTAATCGTATCGAGTTAAAAGGTATTGATTTTAAAGAAGATGGAAACATTCTCGGACACAAACTCGAGTACAACTATAACTCACACAATGTATACATCACGGCAGACAAACAAAAGAATGGAATCAAAGCTAACTTCAAAATTCGCCACAACATTGAAGATGGATCCGTTCAACTAGCAGACCATTATCAACAAAATACTCCAATTGGCGATGGCCCTGTCCTTTTACCAGACAACCATTACCTGTCGACACAATCTGCCCTTTCGAAAGATCCCAACGAAAAGCGTGACCACATGGTCCTTCTTGAGTTTGTAACTGCTGCTGGGATTACACATGGCATGGATGAGCTCTACAAATAA"
	gfpTranslation := "MASKGEELFTGVVPILVELDGDVNGHKFSVSGEGEGDATYGKLTLKFICTTGKLPVPWPTLVTTFSYGVQCFSRYPDHMKRHDFFKSAMPEGYVQERTISFKDDGNYKTRAEVKFEGDTLVNRIELKGIDFKEDGNILGHKLEYNYNSHNVYITADKQKNGIKANFKIRHNIEDGSVQLADHYQQNTPIGDGPVLLPDNHYLSTQSALSKDPNEKRDHMVLLEFVTAAGITHGMDELYK*"

	var writeBuffer bytes.Buffer

	app := application()
	app.Writer = &writeBuffer
	app.Reader = bytes.NewBufferString(gfpDnaSequence)

	args := os.Args[0:1]                   // Name of the program.
	args = append(args, "tr", "-ct", "11") // Append a flag
	err := app.Run(args)

	if err != nil {
		t.Fatalf("Run error: %s", err)
	}
	app.Reader = os.Stdin

	translation := strings.TrimSpace(writeBuffer.String())

	if translation != gfpTranslation {
		t.Errorf("TestTranslationString for string output has failed. Returned %q, want %q", translation, gfpTranslation)
	}

}