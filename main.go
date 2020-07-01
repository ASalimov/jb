package main

import "github.com/ASalimov/jb/cmd"
func main() {
	//s:="привет"
	//fmt.Println(string([]rune(s)[:3]))

	cmd.Execute()
	//n := 20
	//b := bar.NewWithOpts(
	//	bar.WithDimensions(20, 20),
	//	bar.WithFormat(
	//		fmt.Sprintf(
	//			" %sloading...%s :percent :bar %s:rate ops/s%s ",
	//			chalk.Blue,
	//			chalk.Reset,
	//			chalk.Green,
	//			chalk.Reset)))
	//
	//for i := 0; i < n; i++ {
	//	b.Interrupt("asdfasdf")
	//	b.Tick()
	//	time.Sleep(500 * time.Millisecond)
	//}
	//
	//b.Done()
}

