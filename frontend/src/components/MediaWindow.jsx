import { useCallback, useEffect, useRef, useState } from "react";
import {
  getCurrentPlayback,
  getMedia,
} from "../services/api";

function MediaWindow({ windowData, refreshKey }) {
  const [playback, setPlayback] = useState(null);
  const [media, setMedia] = useState([]);
  const videoRef = useRef(null);
  const currentMediaIdRef = useRef(null);

  const loadPlayback = useCallback(async () => {
    try {
      const response = await getCurrentPlayback(windowData.id);
      setPlayback(response.data);
    } catch (error) {
      console.error(
        `Failed to load playback for ${windowData.name}:`,
        error
      );
    }
  }, [windowData.id, windowData.name]);

  useEffect(() => {
    async function load() {
      try {
        const [playbackResponse, mediaResponse] = await Promise.all([
          getCurrentPlayback(windowData.id),
          getMedia(),
        ]);

        setPlayback(playbackResponse.data);
        setMedia(mediaResponse.data);
      } catch (error) {
        console.error("Failed to load window:", error);
      }
    }

    load();
  }, [windowData.id]);

  // Refresh immediately when backend sends an update.
  useEffect(() => {
    if (refreshKey === 0) {
      return;
    }

    loadPlayback();
  }, [refreshKey, loadPlayback]);

  // Keep frontend playback state synchronized with backend.
  useEffect(() => {
    const interval = setInterval(() => {
      loadPlayback();
    }, 1000);

    return () => clearInterval(interval);
  }, [loadPlayback]);

  const currentMedia = media.find(
    (item) => item.id === playback?.MediaID
  );

  // Only seek video when the actual media changes.
  useEffect(() => {
    if (!playback || !currentMedia) {
      return;
    }

    if (currentMediaIdRef.current === playback.mediaId) {
      return;
    }

    currentMediaIdRef.current = playback.mediaId;

    if (
      currentMedia.type === "video" &&
      videoRef.current
    ) {
      videoRef.current.currentTime = playback.offset || 0;

      videoRef.current.play().catch((error) => {
        console.log("Autoplay prevented:", error);
      });
    }
  }, [playback, currentMedia]);

  if (!playback) {
    return (
      <div className="media-window">
        <h3>{windowData.name}</h3>

        <div className="media-content">
          No media
        </div>
      </div>
    );
  }

  if (!currentMedia) {
    return (
      <div className="media-window">
        <h3>{windowData.name}</h3>

        <div className="media-content">
          Media not found
        </div>
      </div>
    );
  }

  return (
    <div className="media-window">
      <h3>{windowData.name}</h3>

      <div className="media-content">
        {currentMedia.type === "image" && (
          <img
            src={currentMedia.url}
            alt={currentMedia.name}
          />
        )}

        {currentMedia.type === "video" && (
          <video
            ref={videoRef}
            src={currentMedia.url}
            muted
            playsInline
            controls={false}
          />
        )}

        {currentMedia.type === "blank" && (
          <div className="blank-media" />
        )}
      </div>

      <div className="media-info">
        <strong>{currentMedia.name}</strong>

        <span>
          {playback.offset}s / {currentMedia.duration}s
        </span>
      </div>
    </div>
  );
}

export default MediaWindow;