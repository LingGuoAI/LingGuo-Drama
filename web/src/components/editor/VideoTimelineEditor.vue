<template>
    <div class="timeline-container">
        <div class="timeline-toolbar">
            <div class="left-tools">
                <t-button shape="square" variant="text" size="small" @click="togglePlay">
                    <template #icon><t-icon :name="isPlaying ? 'pause-circle' : 'play-circle'" size="20px" /></template>
                </t-button>
                <span class="time-code">{{ formatTime(currentTime) }} / {{ formatTime(totalDuration) }}</span>
            </div>
            <div class="right-tools">
                <t-slider v-model="zoomLevel" :min="0.5" :max="4" :step="0.1" style="width: 100px" />
            </div>
        </div>

        <div class="timeline-body" ref="scrollContainer">

            <div class="timeline-ruler" :style="{ width: timelineWidth + 'px' }" @click="handleRulerClick">
                <div v-for="tick in ticks" :key="tick.value" class="tick" :class="tick.type"
                    :style="{ left: tick.left + 'px' }">
                    <span v-if="tick.text" class="tick-text">{{ tick.text }}</span>
                </div>
                <div class="playhead" :style="{ left: playheadX + 'px' }">
                    <div class="playhead-knob"></div>
                    <div class="playhead-line"></div>
                </div>
            </div>

            <div class="tracks-wrapper" :style="{ width: timelineWidth + 'px' }">

                <div class="track-row video-track" @dragover.prevent @drop="onDrop($event, 'video')">
                    <div class="track-header sticky-header"><t-icon name="video" /> Video</div>
                    <div class="track-lane">
                        <div v-for="clip in clips" :key="clip.id" class="clip-item video-clip"
                            :class="{ selected: currentId === clip.id }" :style="getClipStyle(clip)"
                            @click.stop="$emit('select-clip', clip)">
                            <div class="clip-label">{{ clip.name }}</div>

                            <div class="clip-handle left"></div>
                            <div class="clip-handle right"></div>

                            <div class="clip-delete" @click.stop="handleDelete(clip.id)">
                                <t-icon name="close" size="10px" />
                            </div>
                        </div>
                    </div>
                </div>

                <div class="track-row audio-track" @dragover.prevent @drop="onDrop($event, 'audio')">
                    <div class="track-header sticky-header"><t-icon name="sound" /> Audio</div>
                    <div class="track-lane">
                        <div v-for="clip in audioClips" :key="clip.id" class="clip-item audio-clip"
                            :style="getClipStyle(clip)">
                            <div class="clip-label">Audio</div>
                            <div class="clip-delete" @click.stop="handleDelete(clip.id)">
                                <t-icon name="close" size="10px" />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onUnmounted } from 'vue'
import { PlayCircleIcon, PauseCircleIcon, VideoIcon, SoundIcon, CloseIcon } from 'tdesign-icons-vue-next'

const props = defineProps<{
    clips: any[]
    audioClips: any[]
    currentTime: number
    totalDuration: number
    currentId: string | number | null
}>()

const emit = defineEmits(['update:time', 'drop-clip', 'select-clip', 'delete-clip'])

// 状态
const zoomLevel = ref(1)
const basePxPerSec = 20
const isPlaying = ref(false)
const scrollContainer = ref<HTMLElement | null>(null)
let playTimer: any = null

// 计算属性
const pxPerSec = computed(() => basePxPerSec * zoomLevel.value)
const timelineWidth = computed(() => Math.max(props.totalDuration * pxPerSec.value + 200, 1000))
const playheadX = computed(() => props.currentTime * pxPerSec.value)

// 标尺刻度
const ticks = computed(() => {
    const res = []
    const duration = props.totalDuration + 10
    const step = zoomLevel.value < 1 ? 5 : 1
    for (let i = 0; i <= duration; i += step) {
        res.push({
            value: i,
            left: i * pxPerSec.value,
            type: i % 5 === 0 ? 'major' : 'minor',
            text: i % 5 === 0 ? formatTimeSimple(i) : ''
        })
    }
    return res
})

// 交互方法
const getClipStyle = (clip: any) => ({
    left: `${clip.start * pxPerSec.value}px`,
    width: `${clip.duration * pxPerSec.value}px`
})

const handleRulerClick = (e: MouseEvent) => {
    // 简单计算：点击位置 / 像素比 = 时间
    // 注意：这里没有减去左侧 header 宽度，实际应用中可能需要减去 100px
    // 如果 timeline-ruler 本身有 margin-left: 100px，则 e.offsetX 是准确的
    const time = e.offsetX / pxPerSec.value
    emit('update:time', Math.max(0, time))
}

const onDrop = (e: DragEvent, trackType: 'video' | 'audio') => {
    e.preventDefault()
    const raw = e.dataTransfer?.getData('application/json')
    if (!raw) return

    const data = JSON.parse(raw)
    // 计算放置时间
    // 修正计算：需要考虑滚动条位置
    const rect = (e.currentTarget as HTMLElement).getBoundingClientRect()
    // e.clientX 是鼠标相对视口的X
    // rect.left 是轨道元素左边缘
    // 减去 100 (Header宽度)
    const relativeX = e.clientX - rect.left - 100
    const scrollLeft = scrollContainer.value?.scrollLeft || 0

    // 实际 X = 相对轨道的点击位置 + 滚动条位置
    const startTime = Math.max(0, (relativeX + scrollLeft) / pxPerSec.value)

    emit('drop-clip', {
        ...data,
        startTime,
        trackType
    })
}

// 🟢 删除处理
const handleDelete = (clipId: string) => {
    console.log('Delete clip:', clipId)
    emit('delete-clip', clipId)
}

// 播放控制
const togglePlay = () => { isPlaying.value ? pause() : play() }

const play = () => {
    isPlaying.value = true
    const startTs = Date.now()
    const startVideoTime = props.currentTime

    playTimer = setInterval(() => {
        const diff = (Date.now() - startTs) / 1000
        let nextTime = startVideoTime + diff
        if (nextTime >= props.totalDuration) {
            nextTime = props.totalDuration
            pause()
        }
        emit('update:time', nextTime)
    }, 50)
}

const pause = () => {
    isPlaying.value = false
    clearInterval(playTimer)
}

onUnmounted(() => clearInterval(playTimer))

// 工具
const formatTime = (s: number) => {
    const m = Math.floor(s / 60)
    const sec = Math.floor(s % 60)
    const ms = Math.floor((s % 1) * 10)
    return `${m.toString().padStart(2, '0')}:${sec.toString().padStart(2, '0')}.${ms}`
}

const formatTimeSimple = (s: number) => {
    const m = Math.floor(s / 60)
    const sec = s % 60
    return `${m}:${sec.toString().padStart(2, '0')}`
}

// 暴露方法
const scrollToClip = (storyboardId: string | number) => {
    const clip = props.clips.find(c => c.storyboardId === storyboardId)
    if (clip) {
        // 简单跳转，不自动播放
        emit('update:time', clip.start)
    }
}

defineExpose({ scrollToClip })
</script>

<style scoped lang="less">
.timeline-container {
    height: 100%;
    display: flex;
    flex-direction: column;
    background: #1e1e1e;
    color: #ccc;

    .timeline-toolbar {
        height: 40px;
        background: #252525;
        border-bottom: 1px solid #333;
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0 10px;

        .left-tools {
            display: flex;
            align-items: center;
            gap: 10px;
        }

        .time-code {
            font-family: monospace;
            color: #0052d9;
            font-weight: 600;
        }
    }

    .timeline-body {
        flex: 1;
        overflow: auto;
        position: relative;

        .timeline-ruler {
            height: 24px;
            background: #202020;
            border-bottom: 1px solid #333;
            position: relative;
            margin-left: 100px;
            /* 避开轨道头 */
            cursor: pointer;

            .tick {
                position: absolute;
                bottom: 0;
                border-left: 1px solid #555;
                height: 5px;
                font-size: 10px;
                color: #777;
                padding-left: 2px;

                &.major {
                    height: 10px;
                    border-left-color: #888;
                }

                .tick-text {
                    position: absolute;
                    top: -14px;
                    left: 0;
                }
            }

            .playhead {
                position: absolute;
                top: 0;
                bottom: -1000px;
                z-index: 100;
                pointer-events: none;

                .playhead-knob {
                    position: absolute;
                    top: 0;
                    left: -5px;
                    width: 0;
                    height: 0;
                    border-left: 5px solid transparent;
                    border-right: 5px solid transparent;
                    border-top: 10px solid #d10f0f;
                }

                .playhead-line {
                    width: 1px;
                    height: 100%;
                    background: #d10f0f;
                }
            }
        }

        .tracks-wrapper {
            position: relative;

            .track-row {
                height: 80px;
                border-bottom: 1px solid #2a2a2a;
                display: flex;
                position: relative;

                .sticky-header {
                    width: 100px;
                    position: sticky;
                    left: 0;
                    z-index: 20;
                    background: #252525;
                    border-right: 1px solid #333;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    gap: 5px;
                    font-size: 12px;
                    color: #999;
                }

                .track-lane {
                    flex: 1;
                    position: relative;
                    background: rgba(255, 255, 255, 0.02);

                    .clip-item {
                        position: absolute;
                        top: 4px;
                        bottom: 4px;
                        border-radius: 4px;
                        overflow: hidden;
                        cursor: pointer;
                        display: flex;
                        align-items: center;
                        justify-content: center;
                        font-size: 10px;
                        border: 1px solid rgba(255, 255, 255, 0.2);

                        &.video-clip {
                            background: #3d5afe;
                            color: #fff;
                        }

                        &.audio-clip {
                            background: #00897b;
                            color: #fff;
                            height: 30px;
                            top: 25px;
                        }

                        &.selected {
                            border: 2px solid #fff;
                            box-shadow: 0 0 4px rgba(255, 255, 255, 0.5);
                            z-index: 15;
                        }

                        .clip-label {
                            white-space: nowrap;
                            overflow: hidden;
                            text-overflow: ellipsis;
                            padding: 0 4px;
                            pointer-events: none;
                        }

                        /* 调整手柄 */
                        .clip-handle {
                            position: absolute;
                            top: 0;
                            bottom: 0;
                            width: 8px;
                            cursor: col-resize;
                            z-index: 20;

                            &:hover {
                                background: rgba(255, 255, 255, 0.3);
                            }

                            &.left {
                                left: 0;
                            }

                            &.right {
                                right: 0;
                            }
                        }

                        .clip-delete {
                            position: absolute;
                            top: 2px;
                            right: 2px;
                            width: 16px;
                            height: 16px;
                            background: rgba(0, 0, 0, 0.6);
                            border-radius: 50%;
                            display: none;
                            /* 默认隐藏 */
                            align-items: center;
                            justify-content: center;
                            cursor: pointer;
                            z-index: 30;
                            /* 🟢 关键：确保层级最高 */

                            &:hover {
                                background: #d10f0f;
                            }
                        }

                        /* 悬停显示删除按钮 */
                        &:hover .clip-delete {
                            display: flex;
                        }
                    }
                }
            }
        }
    }
}
</style>