# Subtitle Download Feature - Flow Diagram

## System Flow Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                      User Clicks "Watch"                            │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                             ▼
                    ┌────────────────┐
                    │  Load Stream   │
                    │   Information  │
                    └────────┬───────┘
                             │
                             ▼
                    ┌────────────────┐
                    │ Extract        │
                    │ Subtitle URLs  │
                    └────────┬───────┘
                             │
                ┌────────────┴────────────┐
                │                         │
                ▼                         ▼
      ┌──────────────────┐      ┌──────────────────┐
      │ Subtitles Found? │      │ No Subtitles     │
      │      YES         │      │     Found        │
      └────────┬─────────┘      └────────┬─────────┘
               │                         │
               │                         ▼
               │              ┌─────────────────────┐
               │              │  Show Subtitle      │
               │              │  Download Dialog    │
               │              └──────────┬──────────┘
               │                         │
               │                         ▼
               │              ┌─────────────────────┐
               │              │  Search OpenSubs    │
               │              │  API (async)        │
               │              └──────────┬──────────┘
               │                         │
               │              ┌──────────┴──────────┐
               │              │                     │
               │              ▼                     ▼
               │     ┌─────────────────┐   ┌──────────────┐
               │     │ Results Found   │   │ No Results   │
               │     └────────┬────────┘   └──────┬───────┘
               │              │                    │
               │              ▼                    ▼
               │     ┌─────────────────┐   ┌──────────────┐
               │     │ Display Results │   │ Show Message │
               │     │ in Dialog       │   │ "No Results" │
               │     └────────┬────────┘   └──────────────┘
               │              │                         
               │              ▼                         
               │     ┌─────────────────┐               
               │     │ User Selects:   │               
               │     │ 1. Download     │               
               │     │ 2. Play w/o sub │               
               │     │ 3. Cancel       │               
               │     └────────┬────────┘               
               │              │                         
               │     ┌────────┴────────┬────────────┐
               │     │                 │            │
               │     ▼                 ▼            ▼
               │  ┌────────┐    ┌──────────┐  ┌────────┐
               │  │Download│    │Play w/o  │  │ Cancel │
               │  │Subtitle│    │Subtitles │  │        │
               │  └────┬───┘    └────┬─────┘  └────────┘
               │       │             │
               │       ▼             │
               │  ┌────────┐         │
               │  │ Save   │         │
               │  │ to Temp│         │
               │  └────┬───┘         │
               │       │             │
               │       └──────┬──────┘
               │              │
               └──────────────┤
                              │
                              ▼
                    ┌─────────────────┐
                    │ Launch Player   │
                    │ with Subtitle   │
                    └─────────────────┘
```

## Component Interaction Diagram

```
┌──────────────┐         ┌──────────────┐         ┌──────────────┐
│              │         │              │         │              │
│   app.go     │────────▶│ subtitles    │────────▶│ OpenSubs     │
│   (GUI)      │         │ dialog.go    │         │ API          │
│              │         │              │         │              │
└──────┬───────┘         └──────┬───────┘         └──────────────┘
       │                        │
       │                        │
       │                        ▼
       │                 ┌──────────────┐
       │                 │              │
       │                 │ subtitles/   │
       │                 │ manager.go   │
       │                 │              │
       │                 └──────┬───────┘
       │                        │
       │                        │
       │                        ▼
       │                 ┌──────────────┐
       │                 │              │
       │                 │ subtitles/   │
       │                 │ opensubtitles│
       │                 │    .go       │
       │                 └──────┬───────┘
       │                        │
       │                        │
       └────────────────────────┤
                                │
                                ▼
                         ┌──────────────┐
                         │              │
                         │ player/      │
                         │ launcher.go  │
                         │              │
                         └──────────────┘
```

## Dialog State Machine

```
┌─────────────────────┐
│   Dialog Opened     │
│   (Initial State)   │
└──────────┬──────────┘
           │
           ▼
┌─────────────────────┐
│   Searching...      │
│   (Loading)         │
└──────────┬──────────┘
           │
    ┌──────┴──────┐
    │             │
    ▼             ▼
┌─────────┐   ┌─────────┐
│ Results │   │   No    │
│ Found   │   │ Results │
└────┬────┘   └────┬────┘
     │             │
     ▼             ▼
┌─────────┐   ┌─────────┐
│ Waiting │   │ Waiting │
│   for   │   │   for   │
│  User   │   │  User   │
│ Action  │   │ Action  │
└────┬────┘   └────┬────┘
     │             │
     └──────┬──────┘
            │
    ┌───────┼───────┐
    │       │       │
    ▼       ▼       ▼
┌─────┐ ┌─────┐ ┌─────┐
│Down-│ │Play │ │Canc-│
│load │ │ w/o │ │ el  │
└──┬──┘ └──┬──┘ └──┬──┘
   │       │       │
   │       │       ▼
   │       │   ┌─────┐
   │       │   │Exit │
   │       │   └─────┘
   │       │
   ▼       ▼
┌──────────────┐
│ Downloading  │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│    Close     │
│   & Play     │
└──────────────┘
```

## API Request Flow

```
App                     Manager               OpenSubtitles Client           API
│                          │                          │                       │
│─Search Request──────────▶│                          │                       │
│                          │                          │                       │
│                          │─Get Language Pref────────▶                       │
│                          │                          │                       │
│                          │─Build Query──────────────▶                       │
│                          │                          │                       │
│                          │                          │─HTTP GET─────────────▶│
│                          │                          │                       │
│                          │                          │◀─JSON Response───────│
│                          │                          │                       │
│                          │◀─Parse Results───────────│                       │
│                          │                          │                       │
│◀─SubtitleResult[]───────│                          │                       │
│                          │                          │                       │
│                          │                          │                       │
│─Download Request────────▶│                          │                       │
│                          │                          │                       │
│                          │─Download Subtitle────────▶                       │
│                          │                          │                       │
│                          │                          │─HTTP GET─────────────▶│
│                          │                          │                       │
│                          │                          │◀─File Data───────────│
│                          │                          │                       │
│                          │                          │─Save to Temp Dir──┐   │
│                          │                          │                   │   │
│                          │                          │◀──────────────────┘   │
│                          │                          │                       │
│                          │◀─File Path───────────────│                       │
│                          │                          │                       │
│◀─Subtitle Path──────────│                          │                       │
│                          │                          │                       │
```

## Data Structure Flow

```
StreamInfo
    │
    ├─── SubtitleURLs []string
    │         │
    │         ▼
    │    ┌─────────────┐
    │    │ Empty? YES  │
    │    └──────┬──────┘
    │           │
    │           ▼
    │    ┌─────────────────┐
    │    │ Search API      │
    │    └──────┬──────────┘
    │           │
    │           ▼
    │    ┌────────────────────┐
    │    │ SubtitleResult[]   │
    │    │  - ID              │
    │    │  - Language        │
    │    │  - LanguageName    │
    │    │  - MovieName       │
    │    │  - FileName        │
    │    │  - DownloadURL     │
    │    └──────┬─────────────┘
    │           │
    │           ▼
    │    ┌─────────────────┐
    │    │ User Selects    │
    │    │ One Result      │
    │    └──────┬──────────┘
    │           │
    │           ▼
    │    ┌─────────────────┐
    │    │ Download File   │
    │    └──────┬──────────┘
    │           │
    │           ▼
    ├───┬── Subtitle Path
    │   │   (string)
    │   │
    ▼   ▼
┌────────────────┐
│ Player Launch  │
│   Arguments    │
└────────────────┘
```

## Thread/Goroutine Model

```
Main UI Thread                    Worker Goroutines
      │                                  │
      │──Load Stream──────────▶          │
      │                                  │
      │                        ┌─────────▼──────────┐
      │                        │ Fetch Stream Info  │
      │                        └─────────┬──────────┘
      │                                  │
      │◀──Stream Info───────────────────│
      │                                  │
      │──Show Dialog─────────▶          │
      │                                  │
      │                        ┌─────────▼──────────┐
      │                        │ Search Subtitles   │
      │                        └─────────┬──────────┘
      │                                  │
      │◀──fyne.Do(Update UI)─────────────│
      │                                  │
      │──User Clicks Download──▶         │
      │                                  │
      │                        ┌─────────▼──────────┐
      │                        │ Download Subtitle  │
      │                        └─────────┬──────────┘
      │                                  │
      │◀──fyne.Do(Close Dialog)──────────│
      │                                  │
      │──Launch Player─────────▶         │
      │                                  │
```

## Error Handling Flow

```
┌─────────────────┐
│ API Request     │
└────────┬────────┘
         │
    ┌────┴────┐
    │         │
    ▼         ▼
┌─────┐   ┌─────┐
│ OK  │   │Error│
└──┬──┘   └──┬──┘
   │         │
   │         ▼
   │    ┌────────────────┐
   │    │ Network Error? │
   │    └────┬──────┬────┘
   │         │      │
   │      Yes│      │No
   │         │      │
   │         ▼      ▼
   │    ┌─────┐ ┌─────┐
   │    │Show │ │Show │
   │    │Net  │ │API  │
   │    │Error│ │Error│
   │    └──┬──┘ └──┬──┘
   │       │       │
   │       └───┬───┘
   │           │
   │           ▼
   │    ┌────────────┐
   │    │ User Can:  │
   │    │ - Retry    │
   │    │ - Cancel   │
   │    │ - Play w/o │
   │    └────────────┘
   │
   ▼
┌────────────────┐
│ Continue       │
│ Normal Flow    │
└────────────────┘
```

This visual documentation helps understand the complete flow of the subtitle download feature from user interaction to final playback.

