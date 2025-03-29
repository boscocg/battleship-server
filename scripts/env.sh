#!/bin/bash
# Description: Generates the .env file for the project

function print_help {
    echo "Usage: ./scripts/generate-env.sh [-d | -s] <stage>"
    echo "Arguments:"
    echo "  stage - the stage to generate the .env file for"
    echo "  -d    - diff the .env.<stage> file with whats on GCP Secret Manager"
    echo "  -s    - update the .env.<stage> file on GCP Secret Manager"
}

# Script Arguments
# ----------------
PROJECT="battledak"

REQUIRED_CLI=(
    gcloud
    jq
)

main() {
    local PARAM1=$1
    local PARAM2=$2

    validate_arguments $PARAM1 $PARAM2
    validate_libraries

    # If -d flag is passed, diff the .env.dev file with whats on GCP Secret Manager
    if test $PARAM1 = '-d'; then
        diff_env $PARAM2
    elif test $PARAM1 = '-s'; then
        echo -e "⚠️ WARNING: DIFF the .env.${PARAM2} file with whats on GCP Secret Manager before updating"
        echo -e "   -> run: pnpm env:diff:${PARAM2}"
        put_env $PARAM2
    else
        generate_env $PARAM1
    fi

    exit 0
}

# Helper Functions
# ----------------
function generate_env {
    local STAGE=$1
    echo -e "\nGenerating .env.${STAGE} file"
    gcloud secrets versions access latest --secret="${PROJECT}-${STAGE}" > .env.${STAGE} || exit 1
}

function diff_env {
    local STAGE=$1
    local SECRET=$(gcloud secrets versions access latest --secret="${PROJECT}-${STAGE}" || exit 1)

    diff .env.${STAGE} <(echo -e "${SECRET}")
}

function put_env {
    local STAGE=$1
    echo "Putting .env.${STAGE} file to GCP Secret Manager"
    echo -n "$(cat .env.${STAGE})" | gcloud secrets versions add "${PROJECT}-${STAGE}" --data-file=- || exit 1
}

function validate_arguments {
    local ARG1=${1:-''} # flag or stage
    local ARG2=${2:-''} # stage

    if [[ (($ARG1 != '-d' && $ARG1 != '-s') && $ARG2 != '') || (($ARG1 = '-d' || $ARG1 = '-s') && $ARG2 = '') ]]; then
        echo -e "\033[0;31mError: Invalid arguments\033[0m"
        print_help
        exit 1
    fi
}

function validate_libraries {
    for cli in "${REQUIRED_CLI[@]}"; do
        if ! command -v "$cli" &>/dev/null; then
            echo -e "\033[0;31mError: $cli could not be found\033[0m"
            exit 1
        fi
    done
}

main $1 $2
