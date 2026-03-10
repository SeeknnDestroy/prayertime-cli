#!/usr/bin/env sh

set -eu

syft_version="${SYFT_VERSION:-v1.42.1}"
install_attempts="${SYFT_INSTALL_ATTEMPTS:-4}"
install_dir="${SYFT_INSTALL_DIR:-${RUNNER_TEMP:-/tmp}/syft-bin}"

mkdir -p "${install_dir}"

# Retry the full install because transient GitHub release asset failures have
# been the flakiest part of the release-dry-run workflow.
attempt=1
while [ "${attempt}" -le "${install_attempts}" ]; do
  script_path="$(mktemp)"
  trap 'rm -f "${script_path}"' EXIT INT TERM HUP

  if curl --fail --silent --show-error --location --retry 5 --retry-all-errors \
    --output "${script_path}" https://get.anchore.io/syft && \
    rm -f "${install_dir}/syft" && \
    sh "${script_path}" -b "${install_dir}" "${syft_version}" && \
    "${install_dir}/syft" version
  then
    if [ -n "${GITHUB_PATH:-}" ]; then
      printf '%s\n' "${install_dir}" >> "${GITHUB_PATH}"
    fi
    exit 0
  fi

  rm -f "${script_path}"
  trap - EXIT INT TERM HUP

  if [ "${attempt}" -lt "${install_attempts}" ]; then
    sleep_seconds=$((attempt * 5))
    printf 'Syft install attempt %s/%s failed; retrying in %ss\n' "${attempt}" "${install_attempts}" "${sleep_seconds}" >&2
    sleep "${sleep_seconds}"
  fi

  attempt=$((attempt + 1))
done

printf 'Failed to install Syft %s after %s attempts\n' "${syft_version}" "${install_attempts}" >&2
exit 1
