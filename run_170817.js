#!/usr/bin/env node
"use strict";

// Imports
let child_process = require("child_process");
let path = require("path")

// Declare subprocess command
let dstPath = path.join('###net###', 'rackstation3', 'YANDEX', 'Ton') + path.sep
if (path.sep === "/") {
  dstPath = dstPath.replace("###net###", path.sep + path.join('home', 'admin'))
} else {
  dstPath = dstPath.replace("###net###", path.sep)
}
let globalTemplate = `fflite -i "###srcName###" -map_metadata -1 -map_chapters -1 -filter_complex "###filter_complex###" -map [v] -vcodec libx264 -preset medium -crf 17 -pix_fmt yuv420p -g 0 -map [a] -acodec aac -ab 256k -metadata:s:a:0 language=rus -metadata:s:a:0 "handler=Russian 2.0" -disposition:1 default ###time### "###dstPath######dstName###"`
let aspectMap = {
  "11": { num: 1, den: 1 },
  "43": { num: 16, den: 15 },
  "169": { num: 64, den: 45 },
  "1.33": { num: 16, den: 15 },
  "1.78": { num: 64, den: 45 },
  "1.85": { num: 148, den: 100 },
  "2.35": { num: 188, den: 100 },
}
let commandArray = [
  {
    srcName: 'c:\\\\tools\\sd_1996_Ne_nazyvay_menya_malyshkoy__trailer.mpg',
    dstName: 'temp1_plus3.mp4',
    ____vChannel: '[0:v]',
    _vfPrefilter: '',
    ______vfCrop: '',
    _____vfScale: 'auto',
    _____aFilter: '[0:1]tostereo',
    ________time: '',
    _______shift: '-00:00:04:00, 01:00:00:00, 02:00:00:00',
    template: `${globalTemplate}`
  },
]

let DEBUG = 0


let sarNum = 1.0
let sarDen = 1.0
let resW = -1.0
let resH = -1.0
let fps = 0
let fpsNum = 0
let fpsDen = 0

function parsSAR(fname) {
  let cmd = 'fflite -i "' + fname + '"'
  let output = ''
  try {
    output = child_process.execSync(cmd).toString();
  } catch (e) {
    //console.log('!!!!!\n' + e + '\n')
  }
  //console.log('#####\n' + cmd + '\n')
  // console.log('#####\n' + output + '\n')
  let res = output.match(/.*\d{1,2}:\d{1,2}.*Video:\ .*\ (\d{2,5})x(\d{2,5}).*SAR\ (\d{1,4}):(\d{1,4})\ .*\ (\d+)\.{0,1}(\d*)\ fps/)
  //console.log(res[0])
  if (res.length < 7) {
    console.log(output)
    return
  }
  resW = parseInt(res[1], 10)
  resH = parseInt(res[2], 10)
  sarNum = parseInt(res[3], 10)
  sarDen = parseInt(res[4], 10)
  //console.log('zzzz: ' + res[5] + '.' + res[6])
  fps = parseFloat(res[5] + '.' + res[6])
  console.log('input  SAR: ' + sarNum + ':' + sarDen + ' (' + (sarNum / sarDen) + '), ' + resW + 'x' + resH + ', ' + fps + ' fps')
  if (sarNum / sarDen !== 1.0) {
    if (sarNum / sarDen < 1.2) {
      // 4:3 1.066666666
      sarNum = 16
      sarDen = 15
    } else if (sarNum / sarDen > 1.65) {
      // 2.35:1 1.88
      sarNum = 47
      sarDen = 25
    } else {
      // 16:9 1.42222222
      sarNum = 64
      sarDen = 45
    }
  }
  if (fps >= 23.9 && fps < 24.0) {
    fpsNum = 24000
    fpsNum = 1001
  } else if (fps >= 29.9 && fps < 30.0) {
    fpsNum = 30000
    fpsNum = 1001
  } else {
    fpsNum = fps
    fpsDen = 1
  }
  fps = fpsNum / fpsDen

  //timcode2frames('01:01:01:01')
}

function timcode2frames(t) {
  if (t === undefined) {
    return undefined
  }
  t = t.trim()
  let res = t.match(/^[+-]*(\d{2}):(\d{2}):(\d{2}):(\d{2})$/)
  let negative = false
  if (t[0] === '-') {
    negative = true
  }
  if (res === null || res.length < 5) {
    res = t.match(/^[+-]*(\d+)$/)
    if (res === null || res.length < 2) {
      return undefined
    }
    //console.log("~~~~~~~")
    t = parseInt(res, 10)
  } else {
    t = ((parseInt(res[1], 10) * 60 + parseInt(res[2], 10)) * 60 + parseInt(res[3], 10)) * fpsNum / fpsDen + parseInt(res[4], 10)
  }
  if (negative) {
    t = -t
  }
  //console.log("~~~~~~~" + t + '-' + res)
  return t
}

function getVFString(item) {
  let ret = ''
  let x = ''
  x = item.____vChannel.trim()
  if (x !== '') {
    ret = x
  } else {
    ret = '[0:0]'
  }
  x = item._vfPrefilter.trim()
  if (x !== '') {
    ret += x + ','
  }
  x = item.______vfCrop.trim()
  if (x !== '') {
    ret += 'crop=' + x + ','
  }

  x = item._____vfScale.trim()
  if (aspectMap[x] !== undefined || x === "auto") {
    if (aspectMap[x] !== undefined) {
      sarNum = aspectMap[x].num
      sarDen = aspectMap[x].den
    } else if (x !== "auto") {
      sarNum = 0
      sarDen = 0
    }
    if (sarNum / sarDen != 1) {
      //console.log('!!!!!!!!! ' + sarNum / sarDen + ',' + resDen)
      let tr = item.______vfCrop.trim()
      if (tr !== '') {
        resW = parseInt(tr.split(':')[0], 10)
        resH = parseInt(tr.split(':')[1], 10)
      }
      resW = (Math.round(resW * sarNum / (sarDen * 2)) * 2)
      ret += 'scale=' + resW + ':' + resH + ','
    }
  } else {
    sarNum = parseInt(x.split(':')[0], 10)
    sarDen = parseInt(x.split(':')[1], 10)
    ret += 'scale=' + x + ','
  }
  ret += 'setsar=1/1[v],'
  console.log('output SAR: ' + sarNum + ':' + sarDen + ' (' + (sarNum / sarDen) + '), ' + resW + 'x' + resH + ', ' + fpsNum + ':' + fpsDen + ' fps')


  item._____aFilter = item._____aFilter.trim()
  let prefix = item._____aFilter.split('[')[1].split(']')[0]
  //console.log('\n' + item._____aFilter)
  //console.log('\n' + prefix)

  let rest = item._____aFilter.split(']')
  rest.shift()
  //console.log('\nrest-' + rest)

  rest = rest.join(']')
  //console.log('\n' + rest)
  let stream = prefix.split(':')[0]
  let channels = prefix.split(':')[1]
  //console.log('\nchannels:' + channels)
  while (channels.split('-').length > 1) {
    let left = channels.split('-')[0]
    let right = channels.split('-')[1]
    left = left.split(',').pop()
    right = right.split(',').shift()
    let replaceStr = left + '-' + right

    let min = parseInt(left, 10)
    let max = parseInt(right, 10)
    if (left > max) {
      let x = min
      min = max
      max = x
    }
    let progression = ''
    for (let i = min; i <= max; i++) {
      progression += i
      if (i < max) {
        progression += ','
      }
    }
    progression
    channels = channels.replace(replaceStr, progression)
  }
  //console.log('\nrchannels:' + channels)
  let ffchmap = ''
  for (let ch of channels.split(',')) {
    ffchmap += '[' + stream + ':' + ch + ']'
  }

  let aResample = ''
  let preAudio = ''
  let shiftAudio = ''
  x = item._______shift
  if (x !== '') {
    let vals = x.split(',')
    let offset = timcode2frames(vals[0])
    if (offset === undefined) {
      console.log("\x1b[41;1mWrong 1st(offset) parameter argument for __offset\x1b[0m")
      err = 0 / 0
    }
    //console.log(vals.length)
    if (vals.length === 3) {
      let from = timcode2frames(vals[1])
      let to = timcode2frames(vals[2])
      if (from === undefined) {
        console.log("\x1b[41;1mWrong 2nd(form) parameter argument for __offset\x1b[0m")
        err = 0 / 0
      }
      if (to === undefined) {
        console.log("\x1b[41;1mWrong 3d(to) parameter argument for __offset\x1b[0m")
        err = 0 / 0
      }
      if (from / to !== 1) {
        aResample = 'asetrate=' + 48000 * from / to + ',aresample=48000,'
      }
    } else if (vals.length !== 1) {
      console.log("\x1b[41;1mWrong parameters number for __offset\x1b[0m")
      err = 0 / 0
    }
    //fps = 25
    //let offset = parseInt(x, 10)
    let val = Math.abs(offset)
    let seconds = Math.trunc(val / fps)
    console.log('seconds ' + seconds + ", " + offset + " frames (" + fps + " fps)")
    seconds += (val - seconds * fps) * (1 / fps)
    seconds = seconds.toFixed(3)
    console.log('seconds ' + seconds)
    if (offset < 0) {
      shiftAudio = 'atrim=' + seconds + ',asetpts=N/SR/TB'
    } else {
      preAudio = 'aevalsrc=0:c=stereo:s=48000:d=' + seconds + '[silence],'
      shiftAudio = 'concat=n=2:v=0:a=1'
    }
  }
  console.log('@@@preAudio  : ' + preAudio)
  console.log('@@@shiftAudio: ' + shiftAudio)
  console.log('@@@aResample : ' + aResample)

  let aFilter = ffchmap + 'amerge=inputs=' + channels.split(',').length + ',' + rest
  if (aResample !== '') {
    aFilter += aResample
  }
  if (shiftAudio !== '') {
    if (preAudio !== '') {
      aFilter = preAudio + aFilter + '[a1],' + '[silence][a1]' + shiftAudio + '[a]'
    } else {
      aFilter += ',' + shiftAudio + '[a]'
    }
  } else {
    aFilter += '[a]'
  }
  //let aFilter = ffchmap + 'amerge=inputs=' + channels.split(',').length + ',' + rest + '[a]'
  aFilter = aFilter
    .replace("tostereo", "pan=stereo|FL=0.707107*FL+0.707107*FC+0.707107*SL|FR=0.707107*FR+0.707107*FC+0.707107*SR,")
    .replace("2left", "pan=stereo|FL=FL|FR=FL,")
    .replace("2right", "pan=stereo|FL=FR|FR=FR,")
    .replace(',,', ',')
    .replace(',,', ',')
    .replace(',[a]', '[a]')
    .replace(',[a1]', '[a1]')
    .replace('amerge=inputs=1[a]', 'anull[a]')
    .replace('amerge=inputs=1,', '')

  console.log('' + ret)
  console.log(aFilter)
  ret += aFilter

  return ret
}

if (path.sep === "/") {
  dstPath = dstPath.replace("###net###", path.join('~'))
} else {
  dstPath = dstPath.replace("###net###", path.sep)
}

let command = ''
for (let i in commandArray) {
  try {
    if (DEBUG !== 0) {
      console.log('\x1b[41;1m')
      console.log('################################################################################')
      console.log('########!!! DEBUG MODE !!!######################################################')
      console.log('################################################################################\x1b[0m')
      i = commandArray.length - 1
    }
    // Start the subprocess
    console.log(`\n\x1b[42;1mFILE ${+i + 1} of ${commandArray.length}:\x1b[0m`);
    let item = commandArray[i]
    parsSAR(item.srcName)
    command = item.template
      .replace("###srcName###", item.srcName)
      .replace("###dstName###", item.dstName)
      .replace("###dstPath###", dstPath)
      .replace("###time###", item.________time)
      .replace("###filter_complex###", getVFString(item))
    console.log('---------')
    child_process.execSync(command, { stdio: [0, 1, 2] });
  } catch (e) {
    // Error handling
    if (e.message.indexOf("Command failed") === -1) {
      throw e;
    } else {
      console.error(`\x1b[41m${e.message} \x1b[0m`);
    }
  }
  if (DEBUG !== 0) { break }
}
// fflite -r 25 -i "names.txt" -filter_complex "[0:1][0:2][0:3][0:4][0:5][0:6]amerge=inputs=6[rus],[0:7][0:8]amerge=inputs=2[eng],[0:0]scale=720:576,setsar=64/45,unsharp=3:3:0.3:3:3:0[sd]" -map [rus] @ac640 "_q0_AUDIO_RUS51.ac3" -map [eng] @ac640 "_q0_AUDIO_ENG20.ac3" -map [sd] @crf13 "_q0_SD.mp4"

// fflite -r 25 -i "names.txt" -filter_complex "[0:1][0:2]amerge=inputs=2[rus],[0:7][0:8]amerge=inputs=2[eng],[0:0]scale=720:576,setsar=64/45,unsharp=3:3:0.3:3:3:0[sd]" -map [rus] @ac640 "_q0_AUDIO_RUS20.ac3" -map [eng] @ac640 "_q0_AUDIO_ENG20.ac3" -map [sd] @crf13 "_q0_SD.mp4"