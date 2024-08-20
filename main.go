// searchfileany project main.go
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
searchfileany folder [--searchpathregex= -pr=] [--nameregex= -nr=]  [--namepathregex= -npr=] [--contentregex= -cr=] [--directoryonly -do] [--fileonly(default) -fo] [--dirandfile -df] [--deep=] [--timebegin=2006-01-02T15:04:05Z07:00] [--timeend=2006-01-02T15:04:05Z07:00] [--replacewith= match got #0-99 -rw] [--replacewithnorename] [--newpathreplacewith= match got #0-99 -nprw] [--newpathnomove] [--removeResultFile -rf] [--removeResultDir -rd] `,
		)
		return
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
		} else if strings.HasPrefix(os.Args[i], "-df") {
			directoryonly = true
			fileonly = true
		} else if strings.HasPrefix(os.Args[i], "-dp=") {
			deep2, _ := strconv.ParseInt(os.Args[i][len("-dp="):], 10, 64)
			deep = int(deep2)
		} else if strings.HasPrefix(os.Args[i], "-tb=") {
			timebegin, _ = time.Parse(time.RFC3339, os.Args[i][len("-tb="):])
		} else if strings.HasPrefix(os.Args[i], "-te=") {
			timeend, _ = time.Parse(time.RFC3339, os.Args[i][len("-te="):])
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
	searchfileany(&results, sdpath, sdpath, pathre, filenamere, cttre, directoryonly, fileonly, showdetail, deep, 1, timebegin, timeend, replacewithstr, replacewithnorename, newpathreplacewith, newpathnomove, namerestr_mapa, removeResultFile, removeResultDir)
	// for(int i=0;i<results.size();i++){
	//     fmt.Println(results[i])
	// }
	if showdetail {
		fmt.Println("total:", len(results))
	}

	//return a.exec();
}

func searchfileany(results *[]string, rootdir, curdir string, pathregex, filenameregex, contentregex *regexp.Regexp, directoryonly, fileonly, showdetail bool, maxdeep, curdeep int, timebegin, timeend time.Time, replacewithstr string, replacewithnorename bool, newpathreplacewith string, newpathnomove bool, namerestr_mapa, removeResultFile, removeResultDir bool) int {
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
					if showdetail {
						fmt.Println("child file:", fullpath)
					}
					if contentregex != nil {
						bnamepass := false
						if filenameregex != nil {
							if namerestr_mapa == false {
								fnma := filenameregex.FindAllStringSubmatchIndex(fdls[i].Name(), -1)
								if len(fnma) > 0 {
									if showdetail {
										fmt.Println("file name pass:", fullpath)
									}
									bnamepass = true
								}
							} else {
								fnma := filenameregex.FindAllStringSubmatchIndex(fullpath, -1)
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

										fmt.Println(fullpath)
										*results = append(*results, fullpath)
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
							if namerestr_mapa == false {
								fnma := filenameregex.FindAllStringSubmatchIndex(fdls[i].Name(), -1)
								if len(fnma) > 0 {
									if showdetail {
										fmt.Println("found name pass file:", fullpath)
									}
									fmt.Println(fullpath)
									*results = append(*results, fullpath)

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
								fnma := filenameregex.FindAllStringSubmatchIndex(fullpath, -1)
								if len(fnma) > 0 {
									if showdetail {
										fmt.Println("found name pass file:", fullpath)
									}
									fmt.Println(fullpath)
									*results = append(*results, fullpath)

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
					if showdetail {
						fmt.Println("child folder:", fullpath)
					}
					if !(directoryonly == false && fileonly == true) {
						if filenameregex != nil {
							if namerestr_mapa == false {
								fnma := filenameregex.FindAllStringSubmatchIndex(fdls[i].Name(), -1)
								if len(fnma) > 0 {
									if showdetail {
										fmt.Println("found name pass folder:", fullpath)
									}
									fmt.Println(fullpath)
									*results = append(*results, fullpath)

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
								fnma := filenameregex.FindAllStringSubmatchIndex(fullpath, -1)
								if len(fnma) > 0 {
									if showdetail {
										fmt.Println("found name pass folder:", fullpath)
									}
									fmt.Println(fullpath)
									*results = append(*results, fullpath)

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

					searchfileany(results, rootdir, fullpath, pathregex, filenameregex, contentregex, directoryonly, fileonly, showdetail, maxdeep, curdeep+1, timebegin, timeend, replacewithstr, replacewithnorename, newpathreplacewith, newpathnomove, namerestr_mapa, removeResultFile, removeResultDir)
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
					if showdetail {
						fmt.Println("child file:", fullpath)
					}
					if contentregex != nil {
						bnamepass := false
						if filenameregex != nil {
							if namerestr_mapa == false {
								fnma := filenameregex.FindAllStringSubmatchIndex(fdls[i].Name(), -1)
								if len(fnma) > 0 {
									if showdetail {
										fmt.Println("file name pass:", fullpath)
									}
									bnamepass = true
								}
							} else {
								fnma := filenameregex.FindAllStringSubmatchIndex(fullpath, -1)
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

										newpath := toolfunc.RegexReplace(fullpath, pama, newpathreplacewith)
										if showdetail {
											fmt.Println("new path:", newpath)
										}
										if fullpath != newpath {
											if newpathnomove == false {
												toolfunc.MoveFile(fullpath, newpath)
											}
											fmt.Println(fullpath, " new path:", newpath)
										} else {
											fmt.Println(fullpath)
										}
										*results = append(*results, fullpath)

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
												newpath := toolfunc.RegexReplace(fullpath, pama, newpathreplacewith)
												if showdetail {
													fmt.Println("new path:", newpath)
												}
												if fullpath != newpath {
													if newpathnomove == false {
														toolfunc.MoveFile(fullpath, newpath)
													}
													fmt.Println(fullpath, " new path:", newpath)
												} else {
													fmt.Println(fullpath)
												}
												*results = append(*results, fullpath)

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
							if namerestr_mapa == false {
								fnma := filenameregex.FindAllStringSubmatchIndex(fdls[i].Name(), -1)
								if len(fnma) > 0 {
									if showdetail {
										fmt.Println("found name pass file:", fullpath)
									}
									newpath := toolfunc.RegexReplace(fullpath, pama, newpathreplacewith)
									if showdetail {
										fmt.Println("new path:", newpath)
									}
									if fullpath != newpath {
										if newpathnomove == false {
											toolfunc.MoveFile(fullpath, newpath)
										}
										fmt.Println(fullpath, " new path:", newpath)
									} else {
										fmt.Println(fullpath)
									}
									*results = append(*results, fullpath)

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
								fnma := filenameregex.FindAllStringSubmatchIndex(fullpath, -1)
								if len(fnma) > 0 {
									if showdetail {
										fmt.Println("found name pass file:", fullpath)
									}
									newpath := toolfunc.RegexReplace(fullpath, pama, newpathreplacewith)
									if showdetail {
										fmt.Println("new path:", newpath)
									}
									if fullpath != newpath {
										if newpathnomove == false {
											toolfunc.MoveFile(fullpath, newpath)
										}
										fmt.Println(fullpath, " new path:", newpath)
									} else {
										fmt.Println(fullpath)
									}
									*results = append(*results, fullpath)

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
							newpath := toolfunc.RegexReplace(fullpath, pama, newpathreplacewith)
							if showdetail {
								fmt.Println("new path:", newpath)
							}
							if fullpath != newpath {
								if newpathnomove == false {
									toolfunc.MoveFile(fullpath, newpath)
								}
								fmt.Println(fullpath, " new path:", newpath)
							} else {
								fmt.Println(fullpath)
							}

							*results = append(*results, fullpath)

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
					if showdetail {
						fmt.Println("child folder:", fullpath)
					}
					if !(directoryonly == false && fileonly == true) {
						if filenameregex != nil {
							if namerestr_mapa == false {
								fnma := filenameregex.FindAllStringSubmatchIndex(fdls[i].Name(), -1)
								if len(fnma) > 0 {
									if showdetail {
										fmt.Println("found name pass folder:", fullpath)
									}
									fmt.Println(fullpath)
									*results = append(*results, fullpath)

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
								fnma := filenameregex.FindAllStringSubmatchIndex(fullpath, -1)
								if len(fnma) > 0 {
									if showdetail {
										fmt.Println("found name pass folder:", fullpath)
									}
									fmt.Println(fullpath)
									*results = append(*results, fullpath)

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

					searchfileany(results, rootdir, fullpath, pathregex, filenameregex, contentregex, directoryonly, fileonly, showdetail, maxdeep, curdeep+1, timebegin, timeend, replacewithstr, replacewithnorename, newpathreplacewith, newpathnomove, namerestr_mapa, removeResultFile, removeResultDir)
				}
			}
		}
	}
	return 0
}
