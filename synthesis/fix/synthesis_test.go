package fix

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/TimothyStiles/poly/synthesis/codon"
	"github.com/TimothyStiles/poly/transform"
)

/******************************************************************************

Synthesis Fixer tests begin here

******************************************************************************/

var dataDir string = "../../data/"

func TestCdsWithAlteredCodonTable(t *testing.T) {
	bla := "ATGAAAAAAAAAAGTATTCAACATTTCCGTGTCGCCCTTATTCCCTTTTTTGCGGCATTTTGCCTTCCTGTTTTTGCTCACCCAGAAACGCTGGTGAAAGTAAAAGATGCTGAAGATCAGTTGGGTGCACGAGTGGGTTACATCGAACTGGATCTCAACAGCGGTAAGATCCTTGAGAGTTTTCGCCCCGAAGAACGTTTTCCAATGATGAGCACTTTTAAAGTTCTGCTATGTGGCGCGGTATTATCCCGTATTGACGCCGGGCAAGAGCAACTCGGTCGCCGCATACACTATTCTCAGAATGACTTGGTTGAGTACTCACCAGTCACAGAAAAGCATCTTACGGATGGCATGACAGTAAGAGAATTATGCAGTGCTGCCATAACCATGAGTGATAACACTGCGGCCAACTTACTTCTGACAACGATCGGAGGACCGAAGGAGCTAACCGCTTTTTTGCACAACATGGGGGATCATGTAACTCGCCTTGATCGTTGGGAACCGGAGCTGAATGAAGCCATACCAAACGACGAGCGTGACACCACGATGCCTGTAGCAATGGCAACAACGTTGCGCAAACTATTAACTGGCGAACTACTTACTCTAGCTTCCCGGCAACAATTAATAGACTGGATGGAGGCGGATAAAGTTGCAGGACCACTTCTGCGCTCGGCCCTTCCGGCTGGCTGGTTTATTGCTGATAAATCTGGAGCCGGTGAGCGTGGGTCTCGCGGTATCATTGCAGCACTGGGGCCAGATGGTAAGCCCTCCCGTATCGTAGTTATCTACACGACGGGGAGTCAGGCAACTATGGATGAACGAAATAGACAGATCGCTGAGATAGGTGCCTCACTGATTAAGCATTGGTAA"

	codonTable := codon.ReadCodonJSON(dataDir + "alteredPichiaTable.json")

	var functions []func(string, chan DnaSuggestion, *sync.WaitGroup)
	functions = append(functions, RemoveSequence([]string{"CGTGT"}, "Should change to CGA with the Altered Picha Table, because I choose this to be highest"))

	fixedSeq, changes, _ := Cds(bla, codonTable, functions)
	textChange := fmt.Sprintf("Changed position %d from %s to %s for reason: %s. Complete sequence: %s", changes[0].Position, changes[0].From, changes[0].To, changes[0].Reason, fixedSeq)
	shouldChangeTo := "Changed position 9 from CGT to CGA for reason: Should change to CGA with the Altered Picha Table, because I choose this to be highest. Complete sequence: ATGAAAAAAAAAAGTATTCAACATTTCCGAGTCGCCCTTATTCCCTTTTTTGCGGCATTTTGCCTTCCTGTTTTTGCTCACCCAGAAACGCTGGTGAAAGTAAAAGATGCTGAAGATCAGTTGGGTGCACGAGTGGGTTACATCGAACTGGATCTCAACAGCGGTAAGATCCTTGAGAGTTTTCGCCCCGAAGAACGTTTTCCAATGATGAGCACTTTTAAAGTTCTGCTATGTGGCGCGGTATTATCCCGTATTGACGCCGGGCAAGAGCAACTCGGTCGCCGCATACACTATTCTCAGAATGACTTGGTTGAGTACTCACCAGTCACAGAAAAGCATCTTACGGATGGCATGACAGTAAGAGAATTATGCAGTGCTGCCATAACCATGAGTGATAACACTGCGGCCAACTTACTTCTGACAACGATCGGAGGACCGAAGGAGCTAACCGCTTTTTTGCACAACATGGGGGATCATGTAACTCGCCTTGATCGTTGGGAACCGGAGCTGAATGAAGCCATACCAAACGACGAGCGTGACACCACGATGCCTGTAGCAATGGCAACAACGTTGCGCAAACTATTAACTGGCGAACTACTTACTCTAGCTTCCCGGCAACAATTAATAGACTGGATGGAGGCGGATAAAGTTGCAGGACCACTTCTGCGCTCGGCCCTTCCGGCTGGCTGGTTTATTGCTGATAAATCTGGAGCCGGTGAGCGTGGGTCTCGCGGTATCATTGCAGCACTGGGGCCAGATGGTAAGCCCTCCCGTATCGTAGTTATCTATACGACGGGGAGTCAGGCAACTATGGATGAACGAAATAGACAGATCGCTGAGATAGGTGCCTCACTGATTAAGCATTGGTAA"
	if textChange != shouldChangeTo {
		t.Errorf("%s\nshould be\n%s", textChange, shouldChangeTo)
	}
}

func BenchmarkCds(b *testing.B) {
	phusion := "MGHHHHHHHHHHSSGILDVDYITEEGKPVIRLFKKENGKFKIEHDRTFRPYIYALLRDDSKIEEVKKITGERHGKIVRIVDVEKVEKKFLGKPITVWKLYLEHPQDVPTIREKVREHPAVVDIFEYDIPFAKRYLIDKGLIPMEGEEELKILAFDIETLYHEGEEFGKGPIIMISYADENEAKVITWKNIDLPYVEVVSSEREMIKRFLRIIREKDPDIIVTYNGDSFDFPYLAKRAEKLGIKLTIGRDGSEPKMQRIGDMTAVEVKGRIHFDLYHVITRTINLPTYTLEAVYEAIFGKPKEKVYADEIAKAWESGENLERVAKYSMEDAKATYELGKEFLPMEIQLSRLVGQPLWDVSRSSTGNLVEWFLLRKAYERNEVAPNKPSEEEYQRRLRESYTGGFVKEPEKGLWENIVYLDFRALYPSIIITHNVSPDTLNLEGCKNYDIAPQVGHKFCKDIPGFIPSLLGHLLEERQKIKTKMKETQDPIEKILLDYRQKAIKLLANSFYGYYGYAKARWYCKECAESVTAWGRKYIELVWKELEEKFGFKVLYIDTDGLYATIPGGESEEIKKKALEFVKYINSKLPGLLELEYEGFYKRGFFVTKKRYAVIDEEGKVITRGLEIVRRDWSEIAKETQARVLETILKHGDVEEAVRIVKEVIQKLANYEIPPEKLAIYEQITRPLHEYKAIGPHVAVAKKLAAKGVKIKPGMVIGYIVLRGDGPISNRAILAEEYDPKKHKYDAEYYIENQVLPAVLRILEGFGYRKEDLRYQKTRQVGLTSWLNIKKSGTGGGGATVKFKYKGEEKEVDISKIKKVWRVGKMISFTYDEGGGKTGRGAVSEKDAPKELLQMLEKQKK*"
	codonTable := codon.ReadCodonJSON(dataDir + "pichiaTable.json")
	var functions []func(string, chan DnaSuggestion, *sync.WaitGroup)
	functions = append(functions, RemoveSequence([]string{"GAAGAC", "GGTCTC", "GCGATG", "CGTCTC", "GCTCTTC", "CACCTGC"}, "TypeIIS restriction enzyme site."))
	for i := 0; i < b.N; i++ {
		seq, _ := codon.Optimize(phusion, codonTable)
		optimizedSeq, changes, err := Cds(seq, codonTable, functions)
		if err != nil {
			b.Errorf("Failed to fix phusion with error: %s", err)
		}
		for _, cutSite := range []string{"GAAGAC", "GGTCTC", "GCGATG", "CGTCTC", "GCTCTTC", "CACCTGC"} {
			if strings.Contains(optimizedSeq, cutSite) {
				fmt.Println(changes)
				b.Errorf("phusion" + " contains " + cutSite)
			}
			if strings.Contains(transform.ReverseComplement(optimizedSeq), cutSite) {
				fmt.Println(changes)
				b.Errorf("phusion" + " reverse complement contains " + cutSite)
			}
		}
	}
}

func TestReversion(t *testing.T) {
	// Previously, there was an error where BsmBI could get in a loop
	// It would first change CGA -> AGA, then get stuck changing AGA -> AGA
	codonTable := codon.ReadCodonJSON(dataDir + "pichiaTable.json")
	seq := "GGACGAGACGGC"
	var functions []func(string, chan DnaSuggestion, *sync.WaitGroup)
	functions = append(functions, RemoveSequence([]string{"GGTCTC", "CGTCTC"}, "TypeIIS restriction enzyme site."))
	_, _, err := Cds(seq, codonTable, functions)
	if err != nil {
		t.Errorf("Failed with error: %s", err)
	}
}

func TestCds(t *testing.T) {
	codonTable := codon.ReadCodonJSON(dataDir + "pichiaTable.json")
	phusion := "MGHHHHHHHHHHSSGILDVDYITEEGKPVIRLFKKENGKFKIEHDRTFRPYIYALLRDDSKIEEVKKITGERHGKIVRIVDVEKVEKKFLGKPITVWKLYLEHPQDVPTIREKVREHPAVVDIFEYDIPFAKRYLIDKGLIPMEGEEELKILAFDIETLYHEGEEFGKGPIIMISYADENEAKVITWKNIDLPYVEVVSSEREMIKRFLRIIREKDPDIIVTYNGDSFDFPYLAKRAEKLGIKLTIGRDGSEPKMQRIGDMTAVEVKGRIHFDLYHVITRTINLPTYTLEAVYEAIFGKPKEKVYADEIAKAWESGENLERVAKYSMEDAKATYELGKEFLPMEIQLSRLVGQPLWDVSRSSTGNLVEWFLLRKAYERNEVAPNKPSEEEYQRRLRESYTGGFVKEPEKGLWENIVYLDFRALYPSIIITHNVSPDTLNLEGCKNYDIAPQVGHKFCKDIPGFIPSLLGHLLEERQKIKTKMKETQDPIEKILLDYRQKAIKLLANSFYGYYGYAKARWYCKECAESVTAWGRKYIELVWKELEEKFGFKVLYIDTDGLYATIPGGESEEIKKKALEFVKYINSKLPGLLELEYEGFYKRGFFVTKKRYAVIDEEGKVITRGLEIVRRDWSEIAKETQARVLETILKHGDVEEAVRIVKEVIQKLANYEIPPEKLAIYEQITRPLHEYKAIGPHVAVAKKLAAKGVKIKPGMVIGYIVLRGDGPISNRAILAEEYDPKKHKYDAEYYIENQVLPAVLRILEGFGYRKEDLRYQKTRQVGLTSWLNIKKSGTGGGGATVKFKYKGEEKEVDISKIKKVWRVGKMISFTYDEGGGKTGRGAVSEKDAPKELLQMLEKQKK*"
	var functions []func(string, chan DnaSuggestion, *sync.WaitGroup)
	functions = append(functions, RemoveSequence([]string{"GAAGAC", "GGTCTC", "GCGATG", "CGTCTC", "GCTCTTC", "CACCTGC"}, "TypeIIS restriction enzyme site."))
	seq, _ := codon.Optimize(phusion, codonTable)
	optimizedSeq, _, err := Cds(seq, codonTable, functions)
	if err != nil {
		t.Errorf("Failed with error: %s", err)
	}

	for _, cutSite := range []string{"GAAGAC", "GGTCTC", "GCGATG", "CGTCTC", "GCTCTTC", "CACCTGC"} {
		if strings.Contains(optimizedSeq, cutSite) {
			t.Errorf("phusion" + " contains " + cutSite)
		}
		if strings.Contains(transform.ReverseComplement(optimizedSeq), cutSite) {
			t.Errorf("phusion" + " reverse complement contains " + cutSite)
		}
	}

	// Does this flip back in the history?
	fixedSeq, _, _ := CdsSimple("ATGTATTGA", codonTable, []string{"TAT"})
	if fixedSeq != "ATGTACTGA" {
		t.Errorf("Failed to fix ATGTATTGA -> ATGTACTGA")
	}

	// Repeat checking
	blaWithRepeat := "ATGAGTATTCAACATTTCCGTGTCGCCCTTATTCCCTTTTTTGCGGCATTTTGCCTTCCTGTTTTTGCTCACCCAGAAACGCTGGTGAAAGTAAAAGATGCTGAAGATCAGTTGGGTGCACGAGTGGGTTACATCGAACTGGATCTCAACAGCGGTAAGATCCTTGAGAGTTTTCGCCCCGAAGAACGTTTTCCAATGATGAGCACTTTTAAAGTTCTGCTATGTGGCGCGGTATTATCCCGTATTGACGCCGGGCAAGAGCAACTCGGTCGCCGCATACACTATTCTCAGAATGACTTGGTTGAGTACTCACCAGTCACAGAAAAGCATCTTACGGATGGCATGACAGTAAGAGAATTATGCAGTGCTGCCATAACCATGAGTGATAACACTGCGGCCAACTTACTTCTGACAACGATCGGAGGACCGAAGGAGCTAACCGCTTTTTTGCACAACATGGGGGATCATGTAACTCGCCTTGATCGTTGGGAACCGGAGCTGAATGAAGCCATACCAAACGACGAGCGTGACACCACGATGCCTGTAGCAATGGCAACAACGTTGCGCAAACTATTAACTGGCGAACTACTTACTCTAGCTTCCCGGCAACAATTAATAGACTGGATGGAGGCGGATAAAGTTGCAGGACCACTTCTGCGCTCGGCCCTTCCGGCTGGCTGGTTTATTGCTGATAAATCTGGAGCCGGTGAGCGTGGGTCTCGCGGTATCATTGCAGCACTGGGGCCAGATGGTAAGCCCTCCCGTATCGTAGTTATCTACACGACGGGGAGTCAGGCAACTATGGATGAACGAAATAGACAGATCGCTGAGATAGGTGCCTCACTGATTAAGCATTGGGGTGCCTCACTGATTAAGCATTGGTAA"
	functions = append(functions, RemoveRepeat(20))
	blaWithoutRepeat, _, err := Cds(blaWithRepeat, codonTable, functions)
	if err != nil {
		t.Errorf("Failed to remove repeat with error: %s", err)
	}
	targetBlaWithoutRepeat := "ATGAGTATTCAACATTTCCGTGTCGCCCTTATTCCCTTTTTTGCGGCATTTTGCCTTCCTGTTTTTGCTCACCCAGAAACGCTGGTGAAAGTAAAAGATGCTGAAGATCAGTTGGGTGCACGAGTGGGTTACATCGAACTGGATCTCAACAGCGGTAAGATCCTTGAGAGTTTTCGCCCCGAAGAACGTTTTCCAATGATGAGCACTTTTAAAGTTCTGCTATGTGGCGCGGTATTATCCCGTATTGACGCCGGGCAAGAGCAACTCGGTCGCCGCATACACTATTCTCAGAATGACTTGGTTGAGTACTCACCAGTCACAGAAAAGCATCTTACGGATGGCATGACAGTAAGAGAATTATGCAGTGCTGCCATAACCATGAGTGATAACACTGCGGCCAACTTACTTCTGACAACGATCGGAGGACCGAAGGAGCTAACCGCTTTTTTGCACAACATGGGGGATCATGTAACTCGCCTTGATCGTTGGGAACCGGAGCTGAATGAAGCCATACCAAACGACGAGCGTGACACCACGATGCCTGTAGCAATGGCAACAACGTTGCGCAAACTATTAACTGGCGAACTACTTACTCTAGCTTCCCGGCAACAATTAATAGACTGGATGGAGGCGGATAAAGTTGCAGGACCACTTCTGCGCTCGGCCCTTCCGGCTGGCTGGTTTATTGCTGATAAATCTGGAGCCGGTGAGCGTGGATCTCGCGGTATCATTGCAGCACTGGGGCCAGATGGTAAGCCCTCCCGTATCGTAGTTATCTACACGACGGGGAGTCAGGCAACTATGGATGAACGAAATAGACAGATCGCTGAGATAGGTGCCTCACTGATTAAGCATTGGGGTGCTTCACTGATCAAACACTGGTAA"

	if blaWithoutRepeat != targetBlaWithoutRepeat {
		t.Errorf("Expected blaWithoutRepeat sequence %s, got: %s", targetBlaWithoutRepeat, blaWithoutRepeat)
	}

	// Test low and high GC content
	var gcFunctions []func(string, chan DnaSuggestion, *sync.WaitGroup)
	gcFunctions = append(gcFunctions, GcContentFixer(0.90, 0.10))
	fixedSeq, _, err = Cds("GGGCCC", codonTable, gcFunctions)
	if fixedSeq != "GGGCCA" {
		fmt.Println(err)
		t.Errorf("Failed to fix GGGCCC -> GGGCCA. Got %s", fixedSeq)
	}
	fixedSeq, _, _ = Cds("AAATTT", codonTable, gcFunctions)
	if fixedSeq != "AAGTTT" {
		fmt.Println(err)
		t.Errorf("Failed to fix AAATTT -> AAGTTT. Got %s", fixedSeq)
	}
}

func TestCdsBadInput(t *testing.T) {
	// This block tests a sequence that is not divisible by 3
	codonTable := codon.ReadCodonJSON(dataDir + "pichiaTable.json")
	var functions []func(string, chan DnaSuggestion, *sync.WaitGroup)
	functions = append(functions, RemoveSequence([]string{"GAAGAC", "GGTCTC", "GCGATG", "CGTCTC", "GCTCTTC", "CACCTGC"}, "TypeIIS restriction enzyme site"))
	_, _, err := Cds("AT", codonTable, functions)
	if err == nil {
		t.Errorf("Cds should fail with sequence input that is not divisible by 3")
	}

	// This block tests a sequence that has a bad GC bias
	badGcBiasFunc := func(sequence string, c chan DnaSuggestion, wg *sync.WaitGroup) {
		c <- DnaSuggestion{0, 1, "XY", 1, "this should fail"}
		wg.Done()
	}
	_, _, err = Cds("ATG", codonTable, []func(string, chan DnaSuggestion, *sync.WaitGroup){badGcBiasFunc})
	if err == nil {
		t.Errorf("XY should fail as a valid GC bias")
	}

	// This block tests something with no solution space
	_, _, err = Cds("GGG", codonTable, []func(string, chan DnaSuggestion, *sync.WaitGroup){GcContentFixer(0.10, 0.05)})
	if err == nil {
		t.Errorf("There should be no solution to GGG -> less than .10 gc content.")
	}

	// This block tests that any given suggestion will suggest within the confines of the sequence.
	outOfRangePosition := func(sequence string, c chan DnaSuggestion, wg *sync.WaitGroup) {
		c <- DnaSuggestion{1000000000000000000, 10, "GC", 1, "this should fail"}
		wg.Done()
	}
	_, _, err = Cds("ATG", codonTable, []func(string, chan DnaSuggestion, *sync.WaitGroup){outOfRangePosition})
	if err == nil {
		t.Errorf("outOfRangePosition should fail because it is out of range of the start and end positions of the sequence.")
	}
}

func TestBtgZIComplexFix(t *testing.T) {
	complexGene := "ATGAAGCTGATTATTGGCGCAATGCATGAAGAATTGCAGGATTCCATCGCGTTCTATAAGCTGAATAAGGTGGAAAACGAGAAGTTCACCATTTATAAGAATGAAGAGATCATGTTTTGCATTACCGGTATCGGTCTGGTGAACGCGGCGGCGCAGCTGAGCTACATTCTGTCTAAATATGATATTGACTCCATTATTAACATCGGTACCAGCGGCGGTTGCGACAAAGAGCTGAAACAAGGCGACATCCTGATCATCGACAAGATCTATAACAGCGTGGCGGACGCCACCGCATTCGGCTACGCGTACGGCCAAGTTCCGCGTATGCCGAAGTACTATGAAACCAGCAACAAAGATATTATTAAAACCATCAGCAAGGCGAAAATTAAGAATATCGCGAGCTCCGACATCTTCATCCATTCTACGGAGCAAGTGAAGAACTTCATCAATAAAATTGAGGACAAGATTAGCGTCCTGGATATGGAGTGTTTTGCGTATGCTCAGACGGCTTATTTGTTCGAAAAGGAGTTTTCTGTGATTAAAATCATTAGCGACGTCATCGGCGAAAAGAATACCAACAACGTGCAGTTCAACGACTTTATCAAGATTGCCGGTAAGGAGATTTTGGAGATTCTGAAGAAAATTCTG"
	codonTable := codon.ReadCodonJSON(dataDir + "freqB.json")
	restrictionEnzymes := []string{"GAAGAC", "GGTCTC", "GCGATG", "CGTCTC", "GCTCTTC", "CACCTGC"}
	_, _, err := CdsSimple(complexGene, codonTable, restrictionEnzymes)
	if err != nil {
		t.Errorf("Failed to fix complex gene with error: %s", err)
	}
}

func TestBufferFailure(t *testing.T) {
	complexGene := "ATGAAAAAGCTGCTGGCACTGCTGGTTGTGGTCACCTTGACCACCAATGTTGTCGTGGCGGGCGTGGCCATTGCAAACGCGGATAAGAAGAAGCAAAACGACATCCGTATCCTGCAAAGCAAACTGGAGGCAATCCTGAAGAGCAAAACCGATGCGAAGTGGGACGTTTCCGAATTGCAAAAAAAAGTGGATACCGAATTCGGCGAGGGCGAGATTACCGTTAGCTTCAAAGATTATACCAAAGTTACGAGCATTGCAAAGGCTGAATTTATCTTTAAAGCCAACAACAAAAAATACACCGGTCAGCTGACCCTGACCCAGACCTACGAAGTTAAGGATAACAAAGCTGAGGATATCAGCGTCATTAGCACCCGCCTGACGAGCATCCTGAGCGAGAAACCGCACGATGAATGGACCGTTACGGACCTGCAAACCAAGATTGACAGCGAGTTTGGTAATGGTGAGATTGCGGTTAGCGGTGGTACGTATAGCGATGATAACAACTACACCGGCGAAACGAAGAAAAAGGCCGAATTCACGTTCACGGGCAATGCGACCACCGATCCGGAAAACACCCTGAAATACATTGGCGAGATTACGCTGACGCATACGTACACGAAACAAACCGTGATTAGCAACGCTCAGATTAATACGGTGGTGACGGATCTGGCGAATCACGATAAATTCGATAACAAGGAAGCAGCGAAATCCGCGATCGAAGCAGCGTTCGCATACAAGGAGGCGGCAAGCGATGCGGAGCCGACCGGTATCAAGGGTATTGAAAAAGCAGAAGCCAAGTATAACAAGTCCGTGGAAGATGATAAAAGCAGCTTCACCGTGACGTTGACGCTGAGCACGGGTTATGTCCTGGAGCAAACGACCAACACCGTCGAGGTTACGGTGAACTTGATGAGCCGTACCGATATCTCTACCAACGAAGAGTTTAAACAAGAGCTGACCAGCTTTGTGAATGACGAGGCGCACAAGAACCAAGCGTGGACCAAGGACGGTCTGCAAAGCGCGCTGAATACCAAGTATGGCAGCGAGGAGTTTGACGTTACGGAAGACGATAGCACCGTTACGTACGATAATAGCGAACTGGGCAAGAAGACGGAAAAGTTCGTCATCACCGGCCAAGGCAGCAAGGAGAATAATAAACAGTACCAGGGTGAGTTGAAGGTTACGCACGATTACAAGGTTACCGCCAATATTTCTACCATTAAGAACGAGCTGGAGACCATTCTGAAGGATAAAGACTATGAGGAGAAAACGTGGACGCTGGACGAGCTGCAAAAAGCGGTTGATACCGAGTTCAACAAAGGTCAAATTACCGTGGAGGAAGTTATTCTGCTGAAGGATGACAATAGCAATGTGGTTAAAAATACCAAAGAATGGAAATTCATCGGCAATAGCAATGATGAAAACGAATTTGTTTATACCGGTGATGTGACCCTGCCGCACACCTGGAAAAGCTATAAAGTGTTGGCCTCTGATATCCAAACGGCGGCAGAAGTTGCAATTAATGGCAAGAGCTATGCAAATATCGAAGCGGCGCAAGAAGACATCACGAACGCAGTCCAAGCCATCACGGGTGTTGACTCCGTTATTTACCCGACGGAAACGCCGAAAGACTGGAATGATGAAACCATTAAATTTACGGTTACGTTCAAAGAGAACTACGTGATTGAAGGTAAAAATGATTTCAGCGTCAAAGCCCGCGTCGGTAATAGCTCCCAGAATCTGGCGGATATTATTAAGGCGGACGACCTGAAAATCAGCGCGGCAAAAGGCAATGATGCTAGCGCGGTTAAAACCCAAATTGAAACCGTGCTGACCGCTGCGGGTCTGGTGAATGGTACCGATTATGTGGACTTCACCGTGGCGCGTACCGATGATGAGGCTACCACCAGCGTTGAGATCACCGGCAAGGGTAGCGATAAAGTTGTTGATGGTTCCAAAGTTACCTTCGTTGTCACCTGGTCCACCGATTTTTCTAAAGACTTGGCAGACATTATTAAGGCGGACGACCTGAAAATCAGCGCGGCGAAAGGCAATGATGTGAGCACCGTGAAGACCCAAATTGAAACCGTTCTGACGGCTGCGGGTCTGGTGAACGGCACCGACTATGTTGATTTCACCGTCGCGCGTACCGACGATGAAGCGACCACCAGCGTTGAGATTACCGGTAAAGGTAGCGATAAGGTTGTCGACGGCAGCAAGGTTACGTTTGTTGTTACCTGGAGCACCGACTTTAGCAAGGACCTGGCGGACATCATTAAGGCGGACGACTTGAAGATTTCTGCCGCAAAGGGTAATGACGTCAGCACCGTTAAAACCCAAATCGAGACGGTTTTGACCGCAGCGGGTCTGGTGAATGGTACCGATTATGTGGACTTTACGGTGGCACGCACCGACGACGAGGCGACCACCAGCGTGGAAATTACCGGTAAGGGTAGCGACAAGGTTGTTGACGGTAGCAAAGTTACGTTTGTTGTTACGTGGAGCACCGACTTTAGCAAGGATTTGGCAGACATTATCAAAGCCGACGACCTGAAAATTTCTGCGGCCAAGGGCAACGATGTCAGCACCGTTAAGATCCAGATTGAGACCGTGCTGACCGCGGCGGGCCTGGTCAACGGCACCGATTATGTTGATTTCACCGTTGCACGCACCGATGATGAGGCCACGACCAGCGTGGAGATTACCGGTAAGGGTAGCGACAAAGTGGTGGACGGTAGCAAAGTGACCTTCGTTGTGACGTGGAGCATTGATTTCAGCAAAGATCTGGCGGATATTATTAAAGCAGACGACCTGAAGATCTCCGCGGCCAAAGGTAATGATGTTAGCGCGGTCAAGATCCAGATCGAGACGGTTCTGACCGCGGCCGGCTTGGTCAACGGTACGGATTATGTGGACTTCACCGTGGCTCGTACGGATGACGAGGCAACGACCTCTGTGGAGATCACGGGTAAGGGTTCTGATAAGGTTGTCGACGGCAGCAAAGTGACCTTTGTCGTTACCTGGAGCACCGACTTCTCCAAGGACTTGGCAGATATCATTAAGGCCGATGACCTGAAGATCAGCGCTGCGAAAGGTAACGACGTGAGCGCGGTTAAGACCCAAATTGAGACCGTCCTGACCGCAGCGGGCTTGGTTAACGGCACGGATTATGTGGACTTCACCGTTGCACGTACCGATGATGAAGCGACGACCAGCGTCGAGATTACCGGTAAGGGTTCTGACAAAGTGGTTGACGGTAGCAAAGTGACCTTCGTGGTCACCTGGAGCACCGATTTCAGCAAAGATCTGGCGGACATTATTAAAGCGGACGATCTGAAGATCAGCGCGGCCAAGGGCAACGACGTGAGCACGGTGAAAACGCAGATTGAAACCGTGCTGACCGCGGCAGGCCTGGTTAACGGTACCGACTATGTCGACTTCACGGTTGCTCGCACGGACGACGAAGCCACCACCAGCGTGGAGATCACGGGTAAAGGCAGCGATAAGGTTGTGGACGGTAGCAAAGTGACGTTCGTGGTTACCTGGAGCACCGATTTCAGCAAAGACCTGGCCGACATCATCAAGGCAGACGACCTGAAGATCAGCGCAGCTAAGGGCAATGACGACAGCGCTGTTAAGACGCAGATTGAGACCGTGCTGACCGCAGCAGGCCTGGTCAACGGTACGGATTACGTCGACTTTACGGTTGCGCGCACGGACGATGAGGCGACCACCAGCGTTGAAATCACCGGTAAGGGTAGCGATAAAGTCGTCGACGGCAGCAAAGTCACCTTCGTGGTCACCTGGAGCACCGATTTCTCTAAGTATTTGGCGGATATCATCAAGGCAGACGACTTGAAGATTAGCGCGGCAAAGGGCAATGACGCAAGCGCGGTGAAAATCCAGATCGAAACGGTCCTGACCGCCGCAGGCCTGGTCAACGGTACCGACTACGTCGATTTTACCGTCGCACGCACGGACGACGAGGCAACGACCAGCGTCGAAATTACGGGTAAGGGTAGCGACAAAGTTGTGGATGGTAGCAAAGTGACCTTTGTTGTCACCTGGTCCACCGATTTCAGCAAGGATCTGGCAGACATTATTAAAGCGGATGATCTGAAAATCTCCGCCGCGAAAGGCAACGACGTTAGCACCGTTAAAACCCAGATCGAGACGGTCCTGACCGCAGCCGGCCTGGTCAATGGCACGGACTATGTGGACTTCACCGTTGCCCGTACCGACGATGAGGCCACCACCAGCGTTGAGATCACCGGCAAAGGTAGCGATAAGGTGGTTGATGGTAGCAAGGTCACGTTCGTTGTGACCTGGAGCACCGACAGCGGTAACGGTGAAGAGCCGGAGAGCGAAGCACTGAGCATCTTTAGCTATAGCATCATTAGCGATAAGTATTCTAAC"
	codonTable := codon.ReadCodonJSON(dataDir + "freqB.json")
	restrictionEnzymes := []string{"GAAGAC", "GGTCTC", "GCGATG", "CGTCTC", "GCTCTTC", "CACCTGC"}
	_, _, err := CdsSimple(complexGene, codonTable, restrictionEnzymes)
	if err != nil {
		t.Errorf("Failed to fix complex gene with error: %s", err)
	}
}

func TestTooMuchRepeat(t *testing.T) {
	// While this gene cannot be fixed right now, they should be able to be fixed later.
	// Please contribute if you can do this! This will require improvements to the repeat
	// fixing function.
	complexGene := "ATGAAGAAACTGCTGCAACTGCTGGCTGTGCTGTCCCTGACCGCGAGCGTCCTGACCGGCATCGTTTCTTATGAGAGCATGAAAAAACTGAACAAACCGCCGGCGTATAATAAAATCGATCAAAACGAAATTCAAAAGAAGCTGGAAGAGAGCATCAAAAATAAGAACCTGACCGAAGATGAAGCCATCGCCGAGCTGAATAACAGCCTGAAGAATGTGAGCGGTATTAAAACGGTGGAAGCGAAAATTCTGACGAGCTACGCGTTCGAGGAAAAAACGTTCGAGGTTAAAGTGATGCTGGAAGAGAATTACATCTGGGACGACTTGAGCTTTAACGGTGAATTCACCGTGAGCGCGAAGGTTGGTACCTACGACGTGATCAAGAAGGAGGAAATCCAGACCATGCTGAATGAAAGCATCCAAGGCAAAAACCTGACGGAGGACGAAGCTATTGCCGAGCTGAACAATAGCCTGAAAAACGTGAGCGGCATCAAAACCGTTGAAGCCAAAATCCTGACCAGCTATGCGTTCGAGGAAAAAACGTTCGAGGTTAAAGTGATGCTGGAGGAGAACTATGTTTGGGACGACCTGAGCTTTGAGGGTAAGTTCAACGTGAATATCTCCGTTTCTAAAGTCATCAAGATTGATCAGAATGTTATGGAGAAGAGCTTCAAAAGCGCCATCCTGCAGGAGTACGACGAAAGCGAAGCCAAAAAAGCGATCATTGAAACGTTCAACAAGATTATCAATCCGGATCTGACCACGGAGCCGAAAATTGAGATCAAAAAACTGGGTGAAGTTGAATGGGATAAAGAGCATGAAATCACCATTAAGGTGAGCTTGAACACCCATAATTACGAATGGAAAAGCGAGTTCGACGGTGAATTTAAAATCAAAACCGTTCTGAATAGCACGCTGATGTTCTACAAGATCGACAAAGACGAGAACATCCACAGCAAAGAATTTAAAGGCACGAGCAGCAAAGACTGGGATGAAATTGAGTTCACCGAAATCATTGAGTTCGGTTGGTACAACAATGGTCAAGTTTGCGGTATCTTTTTCGAAGAGGACAATAATGAACCGATCAATATCTTCACCCGCTTCAGCGAAGATATTGTTTATCCGAATAAACTGAACGAGAATATCAAAAGCCTGAATTACCTGTTCTATGCGAATTCCAACTCTGGTGACCATTTGTCCGATATCAAAAAATGGGACACGAGCAATGTTAACAGCATGGAGGGCACCTTTAAACTGACCACGTTCAGCAATATTGACCTGAGCGGCTGGAACGTGTCTAATGTTACCAACATGAATTGGATCTTTGCACAGAGCGATATTGTTGATTTTGGTATCTCTAAGTGGAATACGAGCTCCGTGACCGACATGAGCAACATGTTCTACGGTGCTCAAGCGTTTAATGGTGACATTAGCACCAAGGAGGTCGATCAGAATAACGAGAAATACGTCGCCTGGGATACGAGCAAAGTCACCGACATGAGCAACATGTTTAGCGGTAGCAGCGCCTTCAATGGTGACATCTCCAAGTGGAACACCAGCTCCGTCACCAATATGAGCGGCATGTTTAGCGATACCTACGCGTTTAACGGTGACATCAGCAAGTGGAACACGAGCAGCGTCACCGACATGAGCAACATGTTTAGCCGCGCGAGCGCCTTTAACGGCGATATCAGCACCAAGGAGGTTGATCAGAACAACGAAAAATATGTCGCGTGGGACACGAGCAAAGTCACCGATATGAGCAACATGTTCTATCACACGTACGCCTTTAATGGCGATATTAGCAAATGGAACACGAGCAGCGTCACGAACATGTCTAGCATGTTCTCCGACGCTAGCGCTTTTAATGGTGATATCAGCACGAAAGAGGTTGATCAGAATAATGAGAAATACGTCGCCTGGGATACCAGCAAGGTTACCGACATGAGCAACATGTTTTACCATACCTACGCGTTCAACGGCGACATCAGCAAATGGAACACCAGCAGCGTGACGGATATGAGCAACATGTTCCTGGGTGCGCAAAATTTCAACGGTGACATCTCCACCAAAGAGGTTGACCAAAACAACGAAAAATACGTTGCGTGGGATACGTCCAAAGTCACGAACATGAGCGGTATGTTCAGCGAAGCAGAGGCGTTCAATGGCGATATTTCCAAGTGGAATACGTCCAGCGTTACGGACATGAGCAGCATGTTTAGCGGTGCGCAGGCGTTCAACGGTGACATCAGCACCAAAGAGGTGGAGAAAAATAACGAGAAATATGTTGCTTGGGACACCAGCAAAGTGACGGATATGTCCAGCATGTTTAGCGAGACCTACGCCTTTAATGGTGACATCTCCAAATGGAACACGTCCTCTGTCACGAATATGAGCAATATGTTCAGCGGTGCCCAGGCCTTCAACTGTGACATCTCCACCAAAGAGGTTGAGAAAAATAATGAGAAGTACGTGGCATGGGACACCTCCAAGGTTACGGATATGAGCTCCATGTTTTTCGGCGCACAGGCCTTTAATCAGGATATCAGCAAGTGGAATATTAGCAGCGTGACGAACATGAGCTATATGTTCTATCGCGCGCAAGCTTTCAATGTGGACATCTCCAACTGGGATGTCAAAAACGTGGAGTATTTCGCAAACTTCTACCATCAAGGTGGTAATTGGGCTAAAGAACGTCAACCGAAATTTCCGGAGAACAAC"
	codonTable := codon.ReadCodonJSON(dataDir + "freqB.json")
	restrictionEnzymes := []string{"GAAGAC", "GGTCTC", "GCGATG", "CGTCTC", "GCTCTTC", "CACCTGC"}
	_, _, err := CdsSimple(complexGene, codonTable, restrictionEnzymes)
	if err == nil {
		t.Errorf("Succeeded in fixing complexGene")
	}
}

func TestBadCodonTable(t *testing.T) {
	bla := "ATGAAAAAAAAAAGTATTCAACATTTCCGTGTCGCCCTTATTCCCTTTTTTGCGGCATTTTGCCTTCCTGTTTTTGCTCACCCAGAAACGCTGGTGAAAGTAAAAGATGCTGAAGATCAGTTGGGTGCACGAGTGGGTTACATCGAACTGGATCTCAACAGCGGTAAGATCCTTGAGAGTTTTCGCCCCGAAGAACGTTTTCCAATGATGAGCACTTTTAAAGTTCTGCTATGTGGCGCGGTATTATCCCGTATTGACGCCGGGCAAGAGCAACTCGGTCGCCGCATACACTATTCTCAGAATGACTTGGTTGAGTACTCACCAGTCACAGAAAAGCATCTTACGGATGGCATGACAGTAAGAGAATTATGCAGTGCTGCCATAACCATGAGTGATAACACTGCGGCCAACTTACTTCTGACAACGATCGGAGGACCGAAGGAGCTAACCGCTTTTTTGCACAACATGGGGGATCATGTAACTCGCCTTGATCGTTGGGAACCGGAGCTGAATGAAGCCATACCAAACGACGAGCGTGACACCACGATGCCTGTAGCAATGGCAACAACGTTGCGCAAACTATTAACTGGCGAACTACTTACTCTAGCTTCCCGGCAACAATTAATAGACTGGATGGAGGCGGATAAAGTTGCAGGACCACTTCTGCGCTCGGCCCTTCCGGCTGGCTGGTTTATTGCTGATAAATCTGGAGCCGGTGAGCGTGGGTCTCGCGGTATCATTGCAGCACTGGGGCCAGATGGTAAGCCCTCCCGTATCGTAGTTATCTACACGACGGGGAGTCAGGCAACTATGGATGAACGAAATAGACAGATCGCTGAGATAGGTGCCTCACTGATTAAGCATTGGTAA"

	codonTable := codon.ReadCodonJSON(dataDir + "incompletePichiaTable.json")
	_, _, err := CdsSimple(bla, codonTable, []string{"GGTCTC"})
	if err == nil {
		t.Errorf("TestBadCodonTable should fail with 'incomplete codon table'")
	}
}

func TestPanicIndex(t *testing.T) {
	// lastCodon := codonList[len(codonList)-1] used to fail here with a panic
	// panic: runtime error: index out of range [-1]
	// This was because there is a BbsI at the end of the gene and sometimes a
	// synthesis fix would be requested at 825, out of the sequence range.

	codonTable := codon.ReadCodonJSON(dataDir + "freqB.json")
	gene := "ATGTCCGAAAAGAATCTGGAGCACAACCACGGTATCATCAAGGGTATCGATATTGCAGCGGAGGTGCGTAAAGACTTCCTGGAGTACAGCATGAGCGTCATCGTGAGCCGCGCACTGCCGGACCTGAAAGACGGTCTGAAACCGGTTCACCGTCGTATTATCTATGCGATGAACGACCTGGGTATCACCGCGGATAAGCCGCACAAGAAGAGCGCGCGTATCGTCGGTGAAGTTATTGGCAAGTATCACCCGCATGGTGACACCGCAGTTTATGATAGCATGGTTCGTATGGCGCAAGATTTTAGCTACCGCTATCCGCTGGTTGACGGCCACGGTAACTTTGGTAGCATCGACGGTGATGGCGCGGCGGCCATGCGTTACACCGAGGCGCGTTTGGCAAAAGTGTCCAATTTTATTATCAAGGACATCGATATGAATACCGTGCCGTTCGTGGACAACTACGACGCAAGCGAGCGTGAACCGGCTTACCTGACGGGCTATTTCCCGAATCTGCTGGCAAATGGTGCAATGGGTATCGCGGTCGGTATGGCTACCAGCATCCCGCCGCATAATCTGCGTGAGGTGATCTCCGCGATTCATGCATTTATCGATAACAAAGATATCACCATCGATGAGATCCTGGACAATCATATTATGGGTCCGGATTTCCCGACCGGTGCTCTGATGACCAACGGTATTCGTATGCGTGAGGGTTACAAAACGGGTCGCGGTGCGGTGACCATCCGCGCTAAAGTCGCACTGGAAGAGAATGATCGCCATGCGCGCTTCATCATTACGGAGATTCCGTATCAGACCAACAAGGCGAAACTGATTGAGCGCATCGCAGAGCTGGTCAAGACGAAGCAGCTGGAAGGTATCAGCGACATTCGTGACGAGAGCAATTATGAAGGTATTCGCATCGTTATCGAGCTGCGTCGCGACAGCAATCCGGACGTGGTCCTGAGCAAGCTGTACAAATTTACCAACTTGCAAACCACGTATAGCTTGAACCTGCTGAGCCTGCACAATAATATTCCGGTTCTGCTGGACCTGAAAAGCATCATCAAACACTACGTCGACTTTCAGATCAACGTTATTATCAAGCGCAGCATTTTTGAGAAGGATAAGATCGAAAAACGCTTCCACATCCTGGAAGCGCTGGATACCGCGCTGGACAATATCGACGCGATTGTCAACATTCTGCGTAACAGCCCGGAGAGCAACGAGGCTAAAATGAAGCTGACCGAAGCGTTCGGCTTCGATGAAGAACAAAATAAAGCGATCCTGGATATGCGTCTGCAACGTTTGGTCGGTCTGGAACGTGGCAAAATCCAGCAGGAGATGGCGCAGATCAAAGAGCGTATTGACTACCTGACCCTGCTGATTAGCGATGAAACCGAACAGAACAATGTTCTGAAGAACCAGCTGAGCGAAATTGCTGAGAAATTCGGTGACAACCGTCGCACGGAGACGATTGACGAGGGTCTGATGGATATCGAGGATGAGGAACTGATTCCGGACGTGAAGACGATGATTCTGCTGAGCGACGAAGGCTATATTCGTCGTGTGGATCCGGAGGAATTTCGCGTCCAAAAGCGCGGTGGTCGCGGTGTGAGCGTGAACTCCAGCAATGAGGACCCGATTGTTATCGCGACGATGGGTAAGATGCGTGACTGGGTCCTGTTTTTCACGAATAGCGGTAAGGTCTTCCGCACCAAAGCCTACAACATTCGCCAATACAGCCGTACCGCGCGCGGCCTGCCGATCGTGAATTTTCTGAACGGTCTGACCGCGGGCGACAAGATTACCGCGATTCTGCCGCTGCGTAATGTGAAAGAGAAAATGAATTATTTGACCTTTATTACCGAGAAGGGTATGATTAAGAAAACCGATATTAGCCTGTTTGACAATATCAACAAAAACGGTAAAATCGCGATTAACTTGAAAGAGGACGACCAACTGGTGACCGTTTTCGCGACCACGGGCGAGGATACCATCTTTGTGGCAAACAAGAGCGGTAAAGTTATCCGTATTCAGGAAAACATCGTCCGCCCGTTGTCTCGTACGGCATCTGGTGTGAAAGCGATTAAGTTGGACGAGAACGATGTGGTGGTGGGTGCAGTTACGAGCTTCGGTATTGAGAACATTACGACCATTTCCTCCAAGGGTAGCTTCAAAAAGACGAACATCGATGAGTATCGTATCAGCGGCCGTAATGGTAAAGGCATCAAAGTCATGAATCTGAACGAAAAGACCGGTGATTTCAAAAGCATCATTGCGGCACGCGAAACCGATCTGGTTTGTATTATTAGCACGGACGGCAATCTGATTAAGACCAAAGCGAGCGATATCCCGGTGCTGGGCCGTGCGGCTGCCGGCGTGCGTGGTATTCGCCTGTCCGAGGGTAATAAGATTCAGGCCGTTATGCTGGAGTACCGTAAACACGGTGAAGAGAACCAGGAATTCGAGGAAGAC"
	_, _, _ = CdsSimple(gene, codonTable, []string{"GAAGAC", "GGTCTC", "GCGATG", "CGTCTC", "GCTCTTC", "CACCTGC", "CGTCTC"})

}