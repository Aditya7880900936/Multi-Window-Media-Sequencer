import { useEffect, useState } from "react";
import MediaWindow from "./components/MediaWindow";
import {
  addPlaylistItem,
  deletePlaylistItem,
  getMedia,
  getPlaylist,
  getWindows,
  syncMedia,
} from "./services/api";
import { connectWebSocket } from "./services/websocket";
import "./App.css";

function App() {
  const [windows, setWindows] = useState([]);
  const [media, setMedia] = useState([]);
  const [playlists, setPlaylists] = useState({});
  const [refreshKey, setRefreshKey] = useState(0);

  const [selectedMedia, setSelectedMedia] = useState("");
  const [syncDuration, setSyncDuration] = useState(10);
  const [syncing, setSyncing] = useState(false);

  const [selectedWindow, setSelectedWindow] = useState("");
  const [playlistMedia, setPlaylistMedia] = useState("");

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  async function loadData() {
    try {
      setError("");

      const [windowsResponse, mediaResponse] = await Promise.all([
        getWindows(),
        getMedia(),
      ]);

      const loadedWindows = windowsResponse.data;

      setWindows(loadedWindows);
      setMedia(mediaResponse.data);

      const playlistData = {};

      await Promise.all(
        loadedWindows.map(async (window) => {
          const response = await getPlaylist(window.id);
          playlistData[window.id] = response.data;
        })
      );

      setPlaylists(playlistData);
    } catch (err) {
      console.error(err);
      setError(err.message || "Failed to load application");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadData();

    const socket = connectWebSocket((message) => {
      console.log("WebSocket event:", message);

      if (
        message.type === "PLAYLIST_UPDATED" ||
        message.type === "SYNC_PLAY"
      ) {
        setRefreshKey((value) => value + 1);

        if (message.type === "PLAYLIST_UPDATED") {
          loadData();
        }
      }
    });

    return () => {
      socket.close();
    };
  }, []);

  async function handleAddPlaylist() {
    if (!selectedWindow || !playlistMedia) {
      return;
    }

    const selected = media.find(
      (item) => item.id === Number(playlistMedia)
    );

    if (!selected) {
      return;
    }

    try {
      await addPlaylistItem(Number(selectedWindow), {
        mediaId: selected.id,
        position: (playlists[selectedWindow]?.length || 0) + 1,
        duration: selected.duration,
      });

      await loadData();
    } catch (err) {
      setError(err.message || "Failed to add playlist item");
    }
  }

  async function handleDeletePlaylist(windowId, itemId) {
    try {
      await deletePlaylistItem(windowId, itemId);
      await loadData();
    } catch (err) {
      setError(err.message || "Failed to delete playlist item");
    }
  }

  async function handleSync() {
    if (!selectedMedia || syncDuration <= 0) {
      return;
    }

    try {
      setSyncing(true);
      setError("");

      await syncMedia(
        Number(selectedMedia),
        Number(syncDuration)
      );

      setRefreshKey((value) => value + 1);
    } catch (err) {
      setError(err.message || "Failed to start sync");
    } finally {
      setSyncing(false);
    }
  }

  if (loading) {
    return (
      <div className="app-loading">
        Loading Media Sequencer...
      </div>
    );
  }

  return (
    <div className="app">
      <header className="app-header">
        <div>
          <h1>Media Sequencer</h1>
          <p>
            Multi-window synchronized media playback
          </p>
        </div>
      </header>

      {error && (
        <div className="error-message">
          {error}
        </div>
      )}

      <section className="controls">
        <div className="control-card">
          <h2>Sync Playback</h2>

          <div className="control-row">
            <select
              value={selectedMedia}
              onChange={(event) =>
                setSelectedMedia(event.target.value)
              }
            >
              <option value="">Select media</option>

              {media.map((item) => (
                <option key={item.id} value={item.id}>
                  {item.name}
                </option>
              ))}
            </select>

            <input
              type="number"
              min="1"
              value={syncDuration}
              onChange={(event) =>
                setSyncDuration(event.target.value)
              }
            />

            <span>seconds</span>

            <button
              onClick={handleSync}
              disabled={!selectedMedia || syncing}
            >
              {syncing ? "Syncing..." : "Sync All Windows"}
            </button>
          </div>
        </div>

        <div className="control-card">
          <h2>Add To Playlist</h2>

          <div className="control-row">
            <select
              value={selectedWindow}
              onChange={(event) =>
                setSelectedWindow(event.target.value)
              }
            >
              <option value="">Select window</option>

              {windows.map((window) => (
                <option key={window.id} value={window.id}>
                  {window.name}
                </option>
              ))}
            </select>

            <select
              value={playlistMedia}
              onChange={(event) =>
                setPlaylistMedia(event.target.value)
              }
            >
              <option value="">Select media</option>

              {media.map((item) => (
                <option key={item.id} value={item.id}>
                  {item.name}
                </option>
              ))}
            </select>

            <button
              onClick={handleAddPlaylist}
              disabled={!selectedWindow || !playlistMedia}
            >
              Add Media
            </button>
          </div>
        </div>
      </section>

      <main className="windows-grid">
        {windows.map((windowData) => (
          <MediaWindow
            key={windowData.id}
            windowData={windowData}
            refreshKey={refreshKey}
          />
        ))}
      </main>

      <section className="playlists">
        <h2>Playlists</h2>

        <div className="playlist-grid">
          {windows.map((window) => (
            <div
              className="playlist-card"
              key={window.id}
            >
              <h3>{window.name}</h3>

              {(playlists[window.id] || []).length === 0 ? (
                <p>No media configured.</p>
              ) : (
                <div className="playlist-items">
                  {playlists[window.id].map((item) => (
                    <div
                      className="playlist-item"
                      key={item.id}
                    >
                      <div>
                        <strong>
                          {item.position}.{" "}
                          {item.media?.name || "Unknown media"}
                        </strong>

                        <span>
                          {item.duration}s
                        </span>
                      </div>

                      <button
                        onClick={() =>
                          handleDeletePlaylist(
                            window.id,
                            item.id
                          )
                        }
                      >
                        Remove
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>
          ))}
        </div>
      </section>
    </div>
  );
}

export default App;