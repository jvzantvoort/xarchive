#!/bin/bash
#===============================================================================
#
#         FILE:  build.sh
#
#        USAGE:  build.sh
#
#  DESCRIPTION:  Bash script
#
#      OPTIONS:  ---
# REQUIREMENTS:  ---
#         BUGS:  ---
#        NOTES:  ---
#       AUTHOR:  jvzantvoort (John van Zantvoort), john@vanzantvoort.org
#      COMPANY:  JDC
#      CREATED:  2023-08-04
#
# Copyright (C) 2023 John van Zantvoort
#
#===============================================================================

# variables and constants {{{
C_SCRIPTPATH="$(readlink -f "$0")"
C_SCRIPTNAME="$(basename "$C_SCRIPTPATH" .sh)"
C_PACKAGE="xarchive"

declare -A Options
Options["build"]="build a package"
Options["clean"]="clean build results"
Options["doc"]="update documentation"
Options["fmt"]="format the code"
Options["lint"]="lint the code"
Options["tags"]="generate tags for the code"
Options["update"]="update mod and tags"
Options["help"]="print help"

readonly C_SCRIPTPATH
readonly C_SCRIPTNAME
# }}}
# usage {{{
function usage()
{
    local indent
    indent="${C_SCRIPTNAME//?/ }    "
    options=""
    for key in "${!Options[@]}"
    do
      if [[ -z "${options}" ]]
      then
        options="${key}"
      else
        options="${options}|${key}"
      fi
    done

    printf "%s [%s]\n\n" "${C_SCRIPTNAME}.sh" "${options}"

    for key in "${!Options[@]}"
    do
      kval="${Options[${key}]}"
      printf "%s%-12s %s\n" "${indent}" "${key}" "${kval}"
    done
    printf "\n\n"
    exit
}
# }}}
# git functions {{{
function gitroot() { git rev-parse --show-toplevel; }
function gitversion()
{
    git describe --tags --abbrev=0|sed 's/^.*\([0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*\).*/\1/'
}
function gitrevision() { git rev-parse --short HEAD; }
# }}}
# list functions {{{

function listsubdirs()
{
    local root
    root="$(gitroot)"

    find "${root}" -maxdepth 1 -mindepth 1 -type d \
        -not -name vendor \
        -not -name cmd \
        -not -name ref | \
        while read -r target
    do
        find "$target" -name '*.go' -printf "%h\n"
    done|sort -u
}

function listgofiles()
{
    local root
    root="$(gitroot)"

    listsubdirs | while read -r srcdir
    do
        find "${srcdir}" -type f -name '*.go'
    done | sort
    find "${root}/cmd" -type f -name '*.go'
}

function apidocsdir() { printf "%s/docs/api" "$(gitroot)"; }

function buildapidoc()
{
    local root
    local outputdir

    root="$(gitroot)"
    outputdir="$(apidocsdir)"

    pushd "${root}" >/dev/null 2>&1 || exit 1
    listsubdirs | sed "s,${root}\/,," | while read -r dirn
    do
        godocdown "${dirn}" > "${outputdir}/${dirn}.md"
    done
    popd >/dev/null 2>&1 || exit 2

}

function listcommands()
{
    local root
    root="$(gitroot)"
    find "${root}/cmd" -maxdepth 1 -mindepth 1 -type d -printf "%f\n"
}
# }}}

function build_cross()
{
  local goos="$1"
  local arch="$2"
  local ver="$3"
  local ref="$4"
  local command="$5"
  local outputdir="$6"

  GOOS="${goos}" GOARCH="${arch}" \
    go build \
      -o "${outputdir}/${command}" \
      -ldflags "-X main.version=$ver -X main.revision=$rev" "./cmd/${command}"
  retv="$?"
  return "${retv}"
}

function package_cross()
{
  local goos="$1"; shift
  local ver="$1"; shift
  local ref="$1"; shift
  local arch="$1"; shift
  local args=("$@")
  local tempdir
    local root
    root="$(gitroot)"

  tempdir="$(mktemp -d "build_${goos}_${arch}.XXXXX")"

  printf "build %s %s %s/%s\n" "${goos}" "${arch}" "${ver}" "${ref}"
  for cmnd in "${args[@]}"
  do
    build_cross "${goos}" "${arch}" "${ver}" "${ref}" "${cmnd}" "${tempdir}"
  done
  pushd "${tempdir}" >/dev/null 2>&1 || return 1
  zip "${root}/pkg/${C_PACKAGE}_${ver}_${goos}_${arch}.zip" *
  popd >/dev/null 2>&1
  rm -rf "${tempdir}"
  printf "build %s %s %s/%s, end\n" "${goos}" "${arch}" "${ver}" "${ref}"
}

function gofunc_tags()
{
    local root
    root="$(gitroot)"
    listgofiles | xargs gotags > "${root}/tags"
}

function gofunc_clean()
{
    local root
    root="$(gitroot)"

    listcommands | \
        while read -r target
        do
            [[ -f "${root}/${target}" ]] && rm -vf "${root}/${target}"
        done

    [[ -f "${root}/tags" ]] && rm -vf "${root}/tags"
    [[ -d "${root}/pkg" ]] && rm -rvf "${root}/pkg"
}

function gofunc_lint()
{
    listgofiles | while read -r target
    do
        golangci-lint run --no-config --disable-all --enable=goimports \
            --enable=misspell "${target}"
    done
}

function gofunc_fmt()
{
    listgofiles | while read -r target
    do
        goimports -w "${target}"
    done
}

function gofunc_docs()
{
    buildapidoc
}

function gofunc_mod()
{
  pushd "$(gitroot)" > /dev/null 2>&1 || return 1
  if [[ ! -e "go.mod" ]]
  then
    go mod init || return 2
  fi
  go mod tidy || return 3
  go mod vendor || return 4
  popd > /dev/null 2>&1 || return 5
}

function gofunc_build()
{
    local root
    root="$(gitroot)"
  ver="$(gitversion)"
  rev="$(gitrevision)"

  pushd "${root}" >/dev/null 2>&1 || return 1

  gofunc_mod || return 1

  listcommands | while read -r command
      do
          go build -ldflags "-X main.version=$ver -X main.revision=$rev" "./cmd/${command}"
      done
  popd >/dev/null 2>&1 || return 2
}

function gofunc_cross()
{
    ver="$(gitversion)"
    rev="$(gitrevision)"

    gofunc_clean

    gofunc_mod

    mkdir -p pkg
    for arch in "amd64"
    do
        package_cross "darwin" "${arch}" "${ver}" "${arch}" "$(listcommands)"
    done

    for arch in "amd64" "386" "arm64" "arm"
    do
        package_cross "linux" "${arch}" "${ver}" "${arch}" "$(listcommands)"
    done
}

#------------------------------------------------------------------------------#
#                                    Main                                      #
#------------------------------------------------------------------------------#

ACTION="${1:-help}"

case "${ACTION}" in
    build)  gofunc_build ;;
    cross)  gofunc_cross ;;
    clean)  gofunc_clean ;;
    doc)    gofunc_docs  ;;
    fmt)    gofunc_fmt   ;;
    lint)   gofunc_lint  ;;
    tags)   gofunc_tags  ;;
    update) gofunc_mod
            gofunc_tags  ;;

    help) usage ;;
    *) usage ;;
esac

#------------------------------------------------------------------------------#
#                                  The End                                     #
#------------------------------------------------------------------------------#
# vim: foldmethod=marker
