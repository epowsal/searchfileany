// searchfileany project main.go
//
//auther:iwlb@outlook.com
package main

import (
	"bytespool"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"linbo.ga/toolfunc"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println(`help:
searchfileany folder [--searchpathregex= -pr=] [--nameregex= -nr=]  [--namepathregex= -npr=]
 [--contentregex= -cr=] [--directoryonly -do] [--fileonly(default) -fo] [--dirandfile -df -fd] 
[--deep=] [--timebegin=2006-01-02T15:04:05Z07:00] [--timeend=2006-01-02T15:04:05Z07:00] 
[--replacewith= match got #0-99 -rw] [--replacewithnorename] [--newpathreplacewith= match got #0-99 -nprw] 
[--newpathnomove] [--removeResultFile -rf] [--removeResultDir -rd] [--copyto= -ct=] 
[--moveto= -mt=] [--sizerange=min,max -sr=] 【--countsize -cs】 【--countfile -cf】 【--countdir -cd】
【--countfiledir -cfd -cdf】
copy file or directory:
cp source target
move file or directory
mv source target
`,
		)
		return
	}

	if os.Args[1] == "cp" {
		if len(os.Args) == 4 {
			if os.Args[2][0] != '-' && os.Args[3][0] != '-' {
				if toolfunc.FileExists(os.Args[2]) {
					toolfunc.CopyFile(os.Args[2], os.Args[3])
				} else if toolfunc.DirExists(os.Args[2]) {
					toolfunc.CopyDir(os.Args[2], os.Args[3])
				}
				return
			}
		}
	} else if os.Args[1] == "mv" {
		if len(os.Args) == 4 {
			if os.Args[2][0] != '-' && os.Args[3][0] != '-' {
				if toolfunc.FileExists(os.Args[2]) {
					toolfunc.MoveFile(os.Args[2], os.Args[3])
				} else if toolfunc.DirExists(os.Args[2]) {
					toolfunc.MoveDir(os.Args[2], os.Args[3])

				}
				return
			}
		}
	} else if os.Args[1] == "ls" {
		if len(os.Args) >= 3 {
			if os.Args[2][0] != '-' {
				els, _ := os.ReadDir(os.Args[2])
				s1 := []string{}
				s2 := []string{}
				s3 := []string{}
				for _, e := range els {
					in, ine := e.Info()
					if ine == nil {
						df := "f"
						if in.IsDir() {
							df = "d "
						}
						var ode, linkpath string
						if in.IsDir() == false {
							ode = fmt.Sprintf("%04d", in.Mode())
						}
						if 0 != e.Type()&os.ModeDir {
							ode += "dir|"
						}
						if 0 != e.Type()&os.ModeAppend {
							ode = "append"
						}
						if 0 != e.Type()&os.ModeCharDevice {
							ode = "chardev"
						}
						if 0 != e.Type()&os.ModeDevice {
							ode = "dev"
						}
						if 0 != e.Type()&os.ModeExclusive {
							ode = "exclusive"
						}
						if 0 != e.Type()&os.ModeIrregular {
							ode = "irregular"
						}
						if 0 != e.Type()&os.ModeNamedPipe {
							ode = "namedpipe"
						}
						if 0 != e.Type()&os.ModeSocket {
							ode = "soket"
						}
						if 0 != e.Type()&os.ModeSticky {
							ode = "sticky"
						}
						if 0 != e.Type()&os.ModeSymlink {
							ode = "symlink"
							linkpath, _ = os.Readlink(os.Args[2] + "/" + e.Name())
						}
						if 0 != e.Type()&os.ModeTemporary {
							ode = "temp"
						}
						s1 = append(s1, toolfunc.ToIsoDateTime(in.ModTime())+"  "+toolfunc.HumanSize(uint64(in.Size()))+"  "+df)
						s2 = append(s2, ode)
						s3 = append(s3, e.Name()+"  "+linkpath)
					}
				}
				var re *regexp.Regexp
				if len(os.Args) == 4 {
					re = regexp.MustCompile(os.Args[3])
				}
				s := []byte{}
				maxode := 0
				for i := 0; i < len(s1); i += 1 {
					if strings.HasSuffix(s2[i], "|") {
						s2[i] = s2[i][:len(s2[i])-1]
					}
					if len(s2[i]) > maxode {
						maxode = len(s2[i])
					}
				}
				for i := 0; i < len(s1); i += 1 {
					li := s1[i] + "  " + toolfunc.FixLenWithFillRight(s2[i], maxode, ' ') + "  " + s3[i] + "\n"
					if re != nil {
						if !re.MatchString(li) {
							continue
						}
					}
					s = append(s, []byte(li)...)
				}
				fmt.Println(string(s))
				return
			}
		}
	}

	var directoryonly = false
	var fileonly = true
	var showdetail = false
	var replacewithnorename = false
	var newpathnomove = false
	var pathrestr, filenamerestr, cttrestr, replacewithstr, newpathreplacewith string
	var namerestr_mapa, removeResultFile, removeResultDir bool
	var deep = 1 << 30
	var timebegin time.Time
	var timeend time.Time = time.Unix(toolfunc.MAXINT64, 0)
	var copyto, moveto string
	var sizemin, sizemax int64 = 0, toolfunc.MAXINT64
	var countsize *int64
	var countsize2 int64
	var countfile *int64
	var countfile2 int64
	var countdir *int64
	var countdir2 int64

	for i := 0; i < len(os.Args); i += 1 {
		if strings.HasPrefix(os.Args[i], "--searchpathregex=") {
			pathrestr = os.Args[i][len("--searhpathregex="):]
		} else if strings.HasPrefix(os.Args[i], "-pr=") {
			pathrestr = os.Args[i][len("-pr="):]
		} else if strings.HasPrefix(os.Args[i], "--nameregex=") {
			filenamerestr = os.Args[i][len("--nameregex="):]
		} else if strings.HasPrefix(os.Args[i], "--namepathregex=") {
			filenamerestr = os.Args[i][len("--namepathregex="):]
			namerestr_mapa = true
		} else if strings.HasPrefix(os.Args[i], "-npr=") {
			filenamerestr = os.Args[i][len("-npr="):]
			namerestr_mapa = true
		} else if strings.HasPrefix(os.Args[i], "-nr=") {
			filenamerestr = os.Args[i][len("-nr="):]
		} else if strings.HasPrefix(os.Args[i], "--contentregex=") {
			cttrestr = os.Args[i][len("--contentregex="):]
		} else if strings.HasPrefix(os.Args[i], "--deep=") {
			deep2, _ := strconv.ParseInt(os.Args[i][len("--deep="):], 10, 64)
			deep = int(deep2)
		} else if strings.HasPrefix(os.Args[i], "--timebegin=") {
			timebegin, _ = time.Parse(time.RFC3339, os.Args[i][len("--timebegin="):])
		} else if strings.HasPrefix(os.Args[i], "--timeend=") {
			timeend, _ = time.Parse(time.RFC3339, os.Args[i][len("--timeend="):])
		} else if strings.HasPrefix(os.Args[i], "--replacewithstr=") {
			replacewithstr = os.Args[i][len("--replacewith="):]
		} else if strings.HasPrefix(os.Args[i], "-rw=") {
			replacewithstr = os.Args[i][len("-rw="):]
		} else if strings.HasPrefix(os.Args[i], "--newpathreplacewith=") {
			newpathreplacewith = os.Args[i][len("--newpathreplacewith="):]
		} else if strings.HasPrefix(os.Args[i], "-nprw=") {
			newpathreplacewith = os.Args[i][len("-nprw="):]
		} else if strings.HasPrefix(os.Args[i], "--directoryonly") {
			directoryonly = true
		} else if strings.HasPrefix(os.Args[i], "--newpathnomove") {
			newpathnomove = true
		} else if strings.HasPrefix(os.Args[i], "--replacewithnorename") {
			replacewithnorename = true
		} else if strings.HasPrefix(os.Args[i], "--fileonly") {
			fileonly = true
		} else if strings.HasPrefix(os.Args[i], "--showdetail") {
			showdetail = true
		} else if strings.HasPrefix(os.Args[i], "--dirandfile") {
			directoryonly = true
			fileonly = true
		} else if strings.HasPrefix(os.Args[i], "--removeResultFile") {
			removeResultFile = true
		} else if strings.HasPrefix(os.Args[i], "-df") {
			directoryonly = true
			fileonly = true
		} else if strings.HasPrefix(os.Args[i], "-fd") {
			directoryonly = true
			fileonly = true
		} else if strings.HasPrefix(os.Args[i], "--removeResultDir") {
			removeResultDir = true
		} else if strings.HasPrefix(os.Args[i], "-rf") {
			removeResultFile = true
		} else if strings.HasPrefix(os.Args[i], "-rd") {
			removeResultDir = true
		} else if strings.HasPrefix(os.Args[i], "-cr=") {
			cttrestr = os.Args[i][len("-cr="):]
		} else if strings.HasPrefix(os.Args[i], "-do") {
			directoryonly = true
			fileonly = false
		} else if strings.HasPrefix(os.Args[i], "-fo") {
			directoryonly = false
			fileonly = true
		} else if strings.HasPrefix(os.Args[i], "-sd") {
			showdetail = true
		} else if strings.HasPrefix(os.Args[i], "-dp=") {
			deep2, _ := strconv.ParseInt(os.Args[i][len("-dp="):], 10, 64)
			deep = int(deep2)
		} else if strings.HasPrefix(os.Args[i], "-tb=") {
			timebegin, _ = time.Parse(time.RFC3339, os.Args[i][len("-tb="):])
		} else if strings.HasPrefix(os.Args[i], "-te=") {
			timeend, _ = time.Parse(time.RFC3339, os.Args[i][len("-te="):])
		} else if strings.HasPrefix(os.Args[i], "--copyto=") {
			copyto = toolfunc.StdDir(os.Args[i][len("--copyto="):])
		} else if strings.HasPrefix(os.Args[i], "-ct=") {
			copyto = toolfunc.StdDir(os.Args[i][len("-ct="):])
		} else if strings.HasPrefix(os.Args[i], "--copyto=") {
			moveto = toolfunc.StdDir(os.Args[i][len("--moveto="):])
		} else if strings.HasPrefix(os.Args[i], "-ct=") {
			moveto = toolfunc.StdDir(os.Args[i][len("-mt="):])
		} else if strings.HasPrefix(os.Args[i], "--sizerange=") {
			sls := toolfunc.SplitAny(os.Args[i][len("--sizerange="):], ",-")
			if len(sls) == 1 {
				sizemin = toolfunc.Int64FromStr(sls[0])
				sizemax = toolfunc.Int64FromStr(sls[0])
			} else if len(sls) == 2 {
				sizemin = toolfunc.Int64FromStr(sls[0])
				sizemax = toolfunc.Int64FromStr(sls[1])
			}
		} else if strings.HasPrefix(os.Args[i], "-ct=") {
			sls := toolfunc.SplitAny(os.Args[i][len("-ct="):], ",-")
			if len(sls) == 1 {
				sizemin = toolfunc.Int64FromStr(sls[0])
				sizemax = toolfunc.Int64FromStr(sls[0])
			} else if len(sls) == 2 {
				sizemin = toolfunc.Int64FromStr(sls[0])
				sizemax = toolfunc.Int64FromStr(sls[1])
			}
		} else if strings.HasPrefix(os.Args[i], "--countsize") {
			countsize = &countsize2
		} else if strings.HasPrefix(os.Args[i], "-cs") {
			countsize = &countsize2
		} else if strings.HasPrefix(os.Args[i], "-cdf") || strings.HasPrefix(os.Args[i], "-cfd") {
			countfile = &countfile2
			countdir = &countdir2
		} else if strings.HasPrefix(os.Args[i], "--countfile") {
			countfile = &countfile2
		} else if strings.HasPrefix(os.Args[i], "-cf") {
			countfile = &countfile2
		} else if strings.HasPrefix(os.Args[i], "--countdir") {
			countdir = &countdir2
		} else if strings.HasPrefix(os.Args[i], "-cd") {
			countdir = &countdir2
		} else if strings.HasPrefix(os.Args[i], "--countfiledir") {
			countfile = &countfile2
			countdir = &countdir2
		}
	}
	if deep <= 0 {
		deep = 1
	}

	var results []string
	sdpath := toolfunc.ToAbsolutePath(os.Args[1])

	if showdetail {
		fmt.Println("do search parameter path regex:", pathrestr)
		fmt.Println("name regex:", filenamerestr)
		fmt.Println("name path regex:", namerestr_mapa)
		fmt.Println("content regex:", cttrestr)
		fmt.Println("folder only:", directoryonly)
		fmt.Println("file only:", fileonly)
		fmt.Println("root folder:", sdpath)
		fmt.Println("replace with:", replacewithstr)
		fmt.Println("newpathreplacewith:", newpathreplacewith)
		fmt.Println("replacewithnorename", replacewithnorename)
		fmt.Println("timebegin:", timebegin.Format(time.RFC3339))
		fmt.Println("timeend:", timeend.Format(time.RFC3339))
		fmt.Println("copy to:", copyto)
		fmt.Println("move to:", moveto)
		fmt.Println("size min:", sizemin)
		fmt.Println("size max:", sizemax)
		time.Sleep(2 * time.Second)
	}
	if cttrestr != "" && pathrestr == "" {
		pathrestr = ".*"
	}
	if filenamerestr != "" && pathrestr == "" {
		pathrestr = ".*"
	}
	var pathre, filenamere, cttre *regexp.Regexp
	if pathrestr != "" {
		pathre = regexp.MustCompile(pathrestr)
	}
	if filenamerestr != "" {
		filenamere = regexp.MustCompile(filenamerestr)
	}
	if cttrestr != "" {
		cttre = regexp.MustCompile(cttrestr)
	}
	if showdetail {
		fmt.Println("search folder ", sdpath, " results:")
	}
	if showdetail {
		fmt.Println("timeend", timeend)
	}
	searchfileany(&results, sdpath, sdpath, pathre, filenamere, cttre, directoryonly, fileonly, showdetail, deep, 1, timebegin, timeend, replacewithstr, replacewithnorename, newpathreplacewith, newpathnomove, namerestr_mapa, removeResultFile, removeResultDir, copyto, moveto, sizemin, sizemax, countsize, countfile, countdir)
	// for(int i=0;i<results.size();i++){
	//     fmt.Println(results[i])
	// }
	if showdetail {
		fmt.Println("total:", len(results))
	}

	if countsize != nil {
		fmt.Println("count size:", toolfunc.HumanSize(uint64(*countsize)))
	}

	if countfile != nil {
		fmt.Println("count file:", *countfile)
	}

	if countdir != nil {
		fmt.Println("count dir:", *countdir)
	}

	//return a.exec();
}

func searchfileany(results *[]string, rootdir, curdir string, pathregex, filenameregex, contentregex *regexp.Regexp, directoryonly, fileonly, showdetail bool, maxdeep, curdeep int, timebegin, timeend time.Time, replacewithstr string, replacewithnorename bool, newpathreplacewith string, newpathnomove bool, namerestr_mapa, removeResultFile, removeResultDir bool, copyto, moveto string, sizemin, sizemax int64, countsize, countfile, countdir *int64) int {
	if curdeep > maxdeep {
		if showdetail {
			fmt.Println("curdeep > maxdeep", curdeep, maxdeep)
		}
		return 0
	}
	fdls, _ := os.ReadDir(curdir)
	if showdetail {
		fmt.Println("do search folder:", curdir, " child file size:", len(fdls))
	}
	for i := 0; i < len(fdls); i++ {
		fullpath := curdir + fdls[i].Name()
		if fdls[i].IsDir() {
			fullpath += "/"
		}
		if fdls[i].Name() == "." || fdls[i].Name() == ".." || curdir == fullpath {
			continue
		}
		if pathregex == nil {
			continue
		}
		finfo, _ := fdls[i].Info()
		pama := pathregex.FindAllStringSubmatchIndex(fullpath, -1)
		if newpathreplacewith == "" {
			if len(pama) > 0 {
				if fdls[i].IsDir() == false && !(directoryonly == true && fileonly == false) {
					if timebegin.IsZero() == false && finfo.ModTime().Unix()-timebegin.Unix() < 0 {
						if showdetail {
							fmt.Println("before begin time continue.", finfo.ModTime().Format(time.RFC3339))
						}
						continue
					}
					if timeend.IsZero() && finfo.ModTime().Unix()-timeend.Unix() > 0 {
						if showdetail {
							fmt.Println("after end time continue", finfo.ModTime().Format(time.RFC3339))
						}
						continue
					}
					if !(finfo.Size() >= sizemin && finfo.Size() <= sizemax) {
						continue
					}
					// if showdetail {
					// 	fmt.Println("child file:", fullpath)
					// }
					var fnma [][]int
					if contentregex != nil {
						bnamepass := false
						if filenameregex != nil {
							if namerestr_mapa == false {
								fnma = filenameregex.FindAllStringSubmatchIndex(fdls[i].Name(), -1)
								if len(fnma) > 0 {
									if showdetail {
										fmt.Println("file name pass:", fullpath)
									}
									bnamepass = true
								}
							} else {
								fnma = filenameregex.FindAllStringSubmatchIndex(fullpath, -1)
								if len(fnma) > 0 {
									if showdetail {
										fmt.Println("file name pass:", fullpath)
									}
									bnamepass = true
								}
							}
						} else {
							bnamepass = true
						}
						if bnamepass {
							ff, ffe := os.OpenFile(fullpath, os.O_RDONLY, 0666)
							if ffe != nil {
								continue
							}
							if ff != nil {
								ffsize, _ := ff.Seek(0, os.SEEK_END)
								ff.Seek(0, os.SEEK_SET)
								var prerd []byte
								bclose := false
								curd2 := bytespool.Get(32 * 1024 * 1024)
								defer bytespool.Put(curd2)
								for fi := int64(0); fi < ffsize; fi += 32 * 1024 * 1024 {
									rdn, _ := ff.Read(curd2[1024 : 32*1024*1024])
									if rdn <= 0 {
										break
									}
									var curd []byte
									if len(prerd) > 1024 {
										copy(curd2[:1024], prerd[len(prerd)-1024:])
										curd = curd2[:1024+rdn]
									} else {
										curd = curd2[1024 : 1024+rdn]
									}
									ffnma := contentregex.FindAllStringSubmatchIndex(string(curd), -1)
									//fmt.Println("len(prerd) == 0 && ffsize == int64(len(curd)) && replacewithstr ", len(prerd), 0, ffsize, int64(len(curd)), replacewithstr)
									if len(prerd) == 0 && ffsize == int64(len(curd)) && replacewithstr != "" {
										newctt := toolfunc.RegexReplace(string(curd), ffnma, replacewithstr)
										if showdetail {
											if len(newctt) < 128 {
											}
											//fmt.Println("file conntent replace new", newctt)
										} else {
											//fmt.Println("file conntent replace new", newctt[:128])
										}
										fmt.Println("newctt != string(curd) ", newctt != string(curd), newctt, " dkfdslfkjl\n", string(curd))
										if newctt != string(curd) {
											if showdetail {
												fmt.Println("newctt != string(curd)")
											}
											ff2, ff2e := os.OpenFile(fullpath+".sfanew", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
											if ff2e == nil {
												_, ff2we := ff2.Write([]byte(newctt))
												if showdetail {
													fmt.Println("ff2we", ff2we)
												}
												ff2.Close()
												if replacewithnorename == false {
													rnrl := os.Rename(fullpath+".sfanew", fullpath)
													if showdetail {
														fmt.Println("rename result:", rnrl)
													}
												}
											}
										}

										newpath := toolfunc.RegexReplace(fullpath, fnma, newpathreplacewith)
										if showdetail {
											fmt.Println("old path:", fullpath)
											fmt.Println("new path:", newpath)
										}
										if fullpath != newpath {
											if newpathnomove == false {
												toolfunc.MoveFile(fullpath, newpath)
											}
											fmt.Println("regex change path", fullpath, " to:", newpath)
										}

										fmt.Println(fullpath)
										*results = append(*results, fullpath)
										if countsize != nil {
											*countsize += finfo.Size()
										}
										if countfile != nil {
											*countfile += 1
										}
										if copyto != "" {
											toolfunc.CopyFile(fullpath, copyto+fullpath[len(rootdir):])
										} else if moveto != "" {
											toolfunc.MoveFile(fullpath, moveto+fullpath[len(rootdir):])
										}
										if removeResultFile {
											rme := os.Remove(fullpath)
											if rme != nil {
												log.Println("remove file error:", fullpath, rme)
											} else {
												fmt.Println("remove file ok:", fullpath)
											}
										}

										prerd = curd
										break
									} else {
										if len(ffnma) > 0 {
											if directoryonly == false || fileonly == true {
												if showdetail {
													fmt.Println("1found file:", fullpath)
												}
												fmt.Println(fullpath)
												*results = append(*results, fullpath)
												if countsize != nil {
													*countsize += finfo.Size()
												}
												if countfile != nil {
													*countfile += 1
												}
												if copyto != "" {
													toolfunc.CopyFile(fullpath, copyto+fullpath[len(rootdir):])
												} else if moveto != "" {
													toolfunc.MoveFile(fullpath, moveto+fullpath[len(rootdir):])
												}

												if removeResultFile {
													rme := os.Remove(fullpath)
													if rme != nil {
														log.Println("remove file error:", fullpath, rme)
													} else {
														fmt.Println("remove file ok:", fullpath)
													}
												}

											}
											bclose = true
											ff.Close()
											break
										}
										prerd = curd
									}
								}
								if bclose == false {
									ff.Close()
								}
							}
						}
					} else {
						if filenameregex != nil {
							var fnma [][]int
							if namerestr_mapa == false {
								fnma = filenameregex.FindAllStringSubmatchIndex(fdls[i].Name(), -1)
								if len(fnma) > 0 {
									if showdetail {
										fmt.Println("found name pass file:", fullpath)
									}
									fmt.Println(fullpath)
									*results = append(*results, fullpath)
									if countsize != nil {
										*countsize += finfo.Size()
									}
									if countfile != nil {
										*countfile += 1
									}
									if copyto != "" {
										toolfunc.CopyFile(fullpath, copyto+fullpath[len(rootdir):])
									} else if moveto != "" {
										toolfunc.MoveFile(fullpath, moveto+fullpath[len(rootdir):])
									}

									if removeResultFile {
										rme := os.Remove(fullpath)
										if rme != nil {
											log.Println("remove file error:", fullpath, rme)
										} else {
											fmt.Println("remove file ok:", fullpath)
										}
									}
								}
							} else {
								fnma = filenameregex.FindAllStringSubmatchIndex(fullpath, -1)
								if len(fnma) > 0 {
									newpath := toolfunc.RegexReplace(fullpath, fnma, newpathreplacewith)
									if showdetail {
										fmt.Println("old path:", fullpath)
										fmt.Println("new path:", newpath)
									}
									if fullpath != newpath {
										if newpathnomove == false {
											toolfunc.MoveFile(fullpath, newpath)
										}
										fmt.Println("regex change path", fullpath, " to:", newpath)
									}

									if showdetail {
										fmt.Println("found name pass file:", fullpath)
									}
									fmt.Println(fullpath)
									*results = append(*results, fullpath)
									if countsize != nil {
										*countsize += finfo.Size()
									}
									if countfile != nil {
										*countfile += 1
									}
									if copyto != "" {
										toolfunc.CopyFile(fullpath, copyto+fullpath[len(rootdir):])
									} else if moveto != "" {
										toolfunc.MoveFile(fullpath, moveto+fullpath[len(rootdir):])
									}

									if removeResultFile {
										rme := os.Remove(fullpath)
										if rme != nil {
											log.Println("remove file error:", fullpath, rme)
										} else {
											fmt.Println("remove file ok:", fullpath)
										}
									}

								}
							}
						} else {
							if showdetail {
								fmt.Println("show detail found file:", fullpath)
							}
							fmt.Println(fullpath)
							*results = append(*results, fullpath)
							if countsize != nil {
								*countsize += finfo.Size()
							}
							if countfile != nil {
								*countfile += 1
							}
							if copyto != "" {
								toolfunc.CopyFile(fullpath, copyto+fullpath[len(rootdir):])
							} else if moveto != "" {
								toolfunc.MoveFile(fullpath, moveto+fullpath[len(rootdir):])
							}

							if removeResultFile {
								rme := os.Remove(fullpath)
								if rme != nil {
									log.Println("remove file error:", fullpath, rme)
								} else {
									fmt.Println("remove file ok:", fullpath)
								}
							}

						}
					}
				} else if fdls[i].IsDir() == true {
					// if showdetail {
					// 	fmt.Println("child folder:", fullpath)
					// }
					if !(directoryonly == false && fileonly == true) {
						if filenameregex != nil {
							var fnma [][]int
							if namerestr_mapa == false {
								fnma = filenameregex.FindAllStringSubmatchIndex(fdls[i].Name(), -1)
								if len(fnma) > 0 {
									if showdetail {
										fmt.Println("found name pass folder:", fullpath)
									}
									fmt.Println(fullpath)
									*results = append(*results, fullpath)
									if countsize != nil {
										*countsize += toolfunc.GetDirSize(fullpath)
									}
									if countdir != nil {
										*countdir += 1
									}
									if copyto != "" {
										toolfunc.CopyDir(fullpath, copyto+fullpath[len(rootdir):])
									} else if moveto != "" {
										toolfunc.MoveDir(fullpath, moveto+fullpath[len(rootdir):])
									}

									if removeResultDir {
										rme := os.RemoveAll(fullpath)
										if rme != nil {
											log.Println("remove dir error:", fullpath, rme)
										} else {
											fmt.Println("remove dir ok:", fullpath)
										}
									}
								}
							} else {
								fnma = filenameregex.FindAllStringSubmatchIndex(fullpath, -1)
								if len(fnma) > 0 {
									newpath := toolfunc.RegexReplace(fullpath, fnma, newpathreplacewith)
									if showdetail {
										fmt.Println("old path:", fullpath)
										fmt.Println("new path:", newpath)
									}
									if fullpath != newpath {
										if newpathnomove == false {
											toolfunc.MoveDir(fullpath, newpath)
										}
										fmt.Println("regex change path", fullpath, " to:", newpath)
									}

									if showdetail {
										fmt.Println("found name pass folder:", fullpath)
									}
									fmt.Println(fullpath)
									*results = append(*results, fullpath)
									if countsize != nil {
										*countsize += toolfunc.GetDirSize(fullpath)
									}
									if countdir != nil {
										*countdir += 1
									}
									if copyto != "" {
										toolfunc.CopyDir(fullpath, copyto+fullpath[len(rootdir):])
									} else if moveto != "" {
										toolfunc.MoveDir(fullpath, moveto+fullpath[len(rootdir):])
									}

									if removeResultDir {
										rme := os.RemoveAll(fullpath)
										if rme != nil {
											log.Println("remove dir error:", fullpath, rme)
										} else {
											fmt.Println("remove dir ok:", fullpath)
										}
									}

								}
							}
						} else {
							if showdetail {
								fmt.Println("found folder:", fullpath)
							}
							fmt.Println(fullpath)
							*results = append(*results, fullpath)
							if countsize != nil {
								*countsize += toolfunc.GetDirSize(fullpath)
							}
							if countdir != nil {
								*countdir += 1
							}
							if copyto != "" {
								toolfunc.CopyDir(fullpath, copyto+fullpath[len(rootdir):])
							} else if moveto != "" {
								toolfunc.MoveDir(fullpath, moveto+fullpath[len(rootdir):])
							}

							if removeResultDir {
								rme := os.RemoveAll(fullpath)
								if rme != nil {
									log.Println("remove dir error:", fullpath, rme)
								} else {
									fmt.Println("remove dir ok:", fullpath)
								}
							}
						}
					}

					searchfileany(results, rootdir, fullpath, pathregex, filenameregex, contentregex, directoryonly, fileonly, showdetail, maxdeep, curdeep+1, timebegin, timeend, replacewithstr, replacewithnorename, newpathreplacewith, newpathnomove, namerestr_mapa, removeResultFile, removeResultDir, copyto, moveto, sizemin, sizemax, countsize, countfile, countdir)
				}
			}
		} else {
			if len(pama) > 0 {
				if fdls[i].IsDir() == false && !(directoryonly == true && fileonly == false) {
					if timebegin.IsZero() == false && finfo.ModTime().Unix()-timebegin.Unix() < 0 {
						if showdetail {
							fmt.Println("before begin time continue", finfo.ModTime().Format(time.RFC3339))
						}
						continue
					}
					if timeend.IsZero() == false && finfo.ModTime().Unix()-timeend.Unix() > 1 {
						if showdetail {
							fmt.Println("after end time continue", finfo.ModTime().Format(time.RFC3339))
						}
						continue
					}
					if !(finfo.Size() >= sizemin && finfo.Size() <= sizemax) {
						continue
					}
					// if showdetail {
					// 	fmt.Println("child file:", fullpath)
					// }
					var fnma [][]int
					if contentregex != nil {
						bnamepass := false
						if filenameregex != nil {
							if namerestr_mapa == false {
								fnma = filenameregex.FindAllStringSubmatchIndex(fdls[i].Name(), -1)
								if len(fnma) > 0 {
									if showdetail {
										fmt.Println("file name pass:", fullpath)
									}
									bnamepass = true
								}
							} else {
								fnma = filenameregex.FindAllStringSubmatchIndex(fullpath, -1)
								if len(fnma) > 0 {
									if showdetail {
										fmt.Println("file name pass:", fullpath)
									}
									bnamepass = true
								}
							}
						} else {
							bnamepass = true
						}
						if showdetail {
							fmt.Println("bnamepass", bnamepass)
						}
						if bnamepass {
							ff, ffe := os.OpenFile(fullpath, os.O_RDONLY, 0666)
							if ffe == nil {
								ffsize, _ := ff.Seek(0, os.SEEK_END)
								ff.Seek(0, os.SEEK_SET)
								var prerd []byte
								bclose := false
								var curd2 = bytespool.Get(32 * 1024 * 1024)
								defer bytespool.Put(curd2)
								for fi := int64(0); fi < ffsize; fi += 32 * 1024 * 1024 {
									rdn, _ := ff.Read(curd2[1024 : 32*1024*1024])
									if rdn <= 0 {
										break
									}
									var curd []byte
									if len(prerd) > 1024 {
										copy(curd2[:1024], prerd[len(prerd)-1024:])
										curd = curd2[:1024+rdn]
									} else {
										curd = curd2[1024 : 1024+rdn]
									}
									ffnma := contentregex.FindAllStringSubmatchIndex(string(curd), -1)
									if len(prerd) == 0 && ffsize == int64(len(curd)) && replacewithstr != "" {
										newctt := toolfunc.RegexReplace(string(curd), ffnma, replacewithstr)
										if showdetail {
											if len(newctt) < 128 {
											}
											//fmt.Println("file conntent replace new", newctt)
										} else {
											//fmt.Println("file conntent replace new", newctt[:128])
										}
										if showdetail {
											fmt.Println("newctt != string(curd)", newctt != string(curd))
										}
										if newctt != string(curd) {
											ff2, ff2e := os.OpenFile(fullpath+".sfanew", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
											if ff2e == nil {
												ff2.Write([]byte(newctt))
												ff2.Close()
												if replacewithnorename == false {
													rnrl := os.Rename(fullpath+".sfanew", fullpath)
													if showdetail {
														fmt.Println("rename result:", rnrl)
													}
												}
											}
										}

										fmt.Println(fullpath)

										newpath := toolfunc.RegexReplace(fullpath, fnma, newpathreplacewith)
										if showdetail {
											fmt.Println("old path:", fullpath)
											fmt.Println("new path:", newpath)
										}
										if fullpath != newpath {
											if newpathnomove == false {
												toolfunc.MoveFile(fullpath, newpath)
											}
											fmt.Println("regex change path", fullpath, " to:", newpath)
										}

										*results = append(*results, fullpath)
										if countsize != nil {
											*countsize += finfo.Size()
										}
										if countfile != nil {
											*countfile += 1
										}
										if copyto != "" {
											toolfunc.CopyFile(fullpath, copyto+fullpath[len(rootdir):])
										} else if moveto != "" {
											toolfunc.MoveFile(fullpath, moveto+fullpath[len(rootdir):])
										}

										if removeResultFile {
											rme := os.Remove(fullpath)
											if rme != nil {
												log.Println("remove file error:", fullpath, rme)
											} else {
												fmt.Println("remove file ok:", fullpath)
											}
										}

										prerd = curd
										break
									} else {
										if len(ffnma) > 0 {
											if directoryonly == false || fileonly == true {
												if showdetail {
													fmt.Println("3found file:", fullpath)
												}
												newpath := toolfunc.RegexReplace(fullpath, fnma, newpathreplacewith)
												if showdetail {
													fmt.Println("old path:", fullpath)
													fmt.Println("new path:", newpath)
												}
												if fullpath != newpath {
													if newpathnomove == false {
														toolfunc.MoveFile(fullpath, newpath)
													}
													fmt.Println("regex change path", fullpath, " to:", newpath)
												}
												*results = append(*results, fullpath)
												if countsize != nil {
													*countsize += finfo.Size()
												}
												if countfile != nil {
													*countfile += 1
												}
												if copyto != "" {
													toolfunc.CopyFile(fullpath, copyto+fullpath[len(rootdir):])
												} else if moveto != "" {
													toolfunc.MoveFile(fullpath, moveto+fullpath[len(rootdir):])
												}

												if removeResultFile {
													rme := os.Remove(fullpath)
													if rme != nil {
														log.Println("remove file error:", fullpath, rme)
													} else {
														fmt.Println("remove file ok:", fullpath)
													}
												}
											}
											bclose = true
											ff.Close()
											break
										}
										prerd = curd
									}
								}
								if bclose == false {
									ff.Close()
								}
							}
						}
					} else {
						if filenameregex != nil {
							var fnma [][]int
							if namerestr_mapa == false {
								fnma = filenameregex.FindAllStringSubmatchIndex(fdls[i].Name(), -1)
								if len(fnma) > 0 {
									if showdetail {
										fmt.Println("found name pass file:", fullpath)
									}
									newpath := toolfunc.RegexReplace(fullpath, fnma, newpathreplacewith)
									if showdetail {
										fmt.Println("old path:", fullpath)
										fmt.Println("new path:", newpath)
									}

									if fullpath != newpath {
										if newpathnomove == false {
											toolfunc.MoveFile(fullpath, newpath)
										}
										fmt.Println("regex change path", fullpath, " to:", newpath)
									}
									*results = append(*results, fullpath)
									if countsize != nil {
										*countsize += finfo.Size()
									}
									if countfile != nil {
										*countfile += 1
									}
									if copyto != "" {
										toolfunc.CopyFile(fullpath, copyto+fullpath[len(rootdir):])
									} else if moveto != "" {
										toolfunc.MoveFile(fullpath, moveto+fullpath[len(rootdir):])
									}

									if removeResultFile {
										rme := os.Remove(fullpath)
										if rme != nil {
											log.Println("remove file error:", fullpath, rme)
										} else {
											fmt.Println("remove file ok:", fullpath)
										}
									}
								}
							} else {
								fnma = filenameregex.FindAllStringSubmatchIndex(fullpath, -1)
								if len(fnma) > 0 {
									if showdetail {
										fmt.Println("found name pass file:", fullpath)
									}
									newpath := toolfunc.RegexReplace(fullpath, fnma, newpathreplacewith)
									if showdetail {
										fmt.Println("old path:", fullpath)
										fmt.Println("new path:", newpath)
									}
									if fullpath != newpath {
										if newpathnomove == false {
											toolfunc.MoveFile(fullpath, newpath)
										}
										fmt.Println("regex change path", fullpath, " to:", newpath)
									}
									*results = append(*results, fullpath)
									if countsize != nil {
										*countsize += finfo.Size()
									}
									if countfile != nil {
										*countfile += 1
									}
									if copyto != "" {
										toolfunc.CopyFile(fullpath, copyto+fullpath[len(rootdir):])
									} else if moveto != "" {
										toolfunc.MoveFile(fullpath, moveto+fullpath[len(rootdir):])
									}

									if removeResultFile {
										rme := os.Remove(fullpath)
										if rme != nil {
											log.Println("remove file error:", fullpath, rme)
										} else {
											fmt.Println("remove file ok:", fullpath)
										}
									}
								}
							}
						} else {
							if showdetail {
								fmt.Println("4found file:", fullpath)
							}
							newpath := toolfunc.RegexReplace(fullpath, fnma, newpathreplacewith)
							if showdetail {
								fmt.Println("old path:", fullpath)
								fmt.Println("new path:", newpath)
							}
							if fullpath != newpath {
								if newpathnomove == false {
									toolfunc.MoveFile(fullpath, newpath)
								}
								fmt.Println("regex change path", fullpath, " to:", newpath)
							}

							*results = append(*results, fullpath)
							if countsize != nil {
								*countsize += finfo.Size()
							}
							if countfile != nil {
								*countfile += 1
							}
							if copyto != "" {
								toolfunc.CopyFile(fullpath, copyto+fullpath[len(rootdir):])
							} else if moveto != "" {
								toolfunc.MoveFile(fullpath, moveto+fullpath[len(rootdir):])
							}

							if removeResultFile {
								rme := os.Remove(fullpath)
								if rme != nil {
									log.Println("remove file error:", fullpath, rme)
								} else {
									fmt.Println("remove file ok:", fullpath)
								}
							}
						}
					}
				} else if fdls[i].IsDir() == true {
					// if showdetail {
					// 	fmt.Println("child folder:", fullpath)
					// }
					if !(directoryonly == false && fileonly == true) {
						if filenameregex != nil {
							var fnma [][]int
							if namerestr_mapa == false {
								fnma = filenameregex.FindAllStringSubmatchIndex(fdls[i].Name(), -1)
								if len(fnma) > 0 {
									if showdetail {
										fmt.Println("found name pass folder:", fullpath)
									}
									fmt.Println(fullpath)
									*results = append(*results, fullpath)
									if countsize != nil {
										*countsize += toolfunc.GetDirSize(fullpath)
									}
									if countdir != nil {
										*countdir += 1
									}
									if copyto != "" {
										toolfunc.CopyDir(fullpath, copyto+fullpath[len(rootdir):])
									} else if moveto != "" {
										toolfunc.MoveDir(fullpath, moveto+fullpath[len(rootdir):])
									}

									if removeResultDir {
										rme := os.Remove(fullpath)
										if rme != nil {
											log.Println("remove dir error:", fullpath, rme)
										} else {
											fmt.Println("remove dir ok:", fullpath)
										}
									}
								}
							} else {
								fnma = filenameregex.FindAllStringSubmatchIndex(fullpath, -1)
								if len(fnma) > 0 {
									newpath := toolfunc.RegexReplace(fullpath, fnma, newpathreplacewith)
									if showdetail {
										fmt.Println("old path:", fullpath)
										fmt.Println("new path:", newpath)
									}
									if fullpath != newpath {
										if newpathnomove == false {
											toolfunc.MoveDir(fullpath, newpath)
										}
										fmt.Println("regex change path", fullpath, " to:", newpath)
									}

									if showdetail {
										fmt.Println("found name pass folder:", fullpath)
									}
									fmt.Println(fullpath)
									*results = append(*results, fullpath)
									if countsize != nil {
										*countsize += toolfunc.GetDirSize(fullpath)
									}
									if countdir != nil {
										*countdir += 1
									}
									if copyto != "" {
										toolfunc.CopyDir(fullpath, copyto+fullpath[len(rootdir):])
									} else if moveto != "" {
										toolfunc.MoveDir(fullpath, moveto+fullpath[len(rootdir):])
									}

									if removeResultDir {
										rme := os.Remove(fullpath)
										if rme != nil {
											log.Println("remove dir error:", fullpath, rme)
										} else {
											fmt.Println("remove dir ok:", fullpath)
										}
									}
								}
							}
						} else {
							if showdetail {
								fmt.Println("found folder:", fullpath)
							}
							fmt.Println(fullpath)
							*results = append(*results, fullpath)
							if countsize != nil {
								*countsize += toolfunc.GetDirSize(fullpath)
							}
							if countdir != nil {
								*countdir += 1
							}
							if copyto != "" {
								toolfunc.CopyDir(fullpath, copyto+fullpath[len(rootdir):])
							} else if moveto != "" {
								toolfunc.MoveDir(fullpath, moveto+fullpath[len(rootdir):])
							}

							if removeResultDir {
								rme := os.Remove(fullpath)
								if rme != nil {
									log.Println("remove dir error:", fullpath, rme)
								} else {
									fmt.Println("remove dir ok:", fullpath)
								}
							}
						}
					}

					searchfileany(results, rootdir, fullpath, pathregex, filenameregex, contentregex, directoryonly, fileonly, showdetail, maxdeep, curdeep+1, timebegin, timeend, replacewithstr, replacewithnorename, newpathreplacewith, newpathnomove, namerestr_mapa, removeResultFile, removeResultDir, copyto, moveto, sizemin, sizemax, countsize, countfile, countdir)
				}
			}
		}
	}
	return 0
}
