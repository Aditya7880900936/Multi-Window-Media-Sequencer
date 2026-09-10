const API_BASE_URL = "http://localhost:8080/api";

async function request(url, options = {}) {
  const response = await fetch(`${API_BASE_URL}${url}`, {
    headers: {
      "Content-Type": "application/json",
      ...options.headers,
    },
    ...options,
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.error || "Request failed");
  }

  return data;
}

export async function getWindows() {
  return request("/windows");
}

export async function getWindow(windowId) {
  return request(`/windows/${windowId}`);
}

export async function getCurrentPlayback(windowId) {
  return request(`/windows/${windowId}/current`);
}

export async function getPlaylist(windowId) {
  return request(`/windows/${windowId}/playlist`);
}

export async function getMedia() {
  return request("/media");
}

export async function addPlaylistItem(windowId, item) {
  return request(`/windows/${windowId}/playlist`, {
    method: "POST",
    body: JSON.stringify(item),
  });
}

export async function deletePlaylistItem(windowId, itemId) {
  return request(`/windows/${windowId}/playlist/${itemId}`, {
    method: "DELETE",
  });
}

export async function syncMedia(mediaId, duration) {
  return request("/sync", {
    method: "POST",
    body: JSON.stringify({
      mediaId,
      duration,
    }),
  });
}