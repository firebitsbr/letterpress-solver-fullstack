import { Component, OnInit } from '@angular/core';
import { Http, Headers } from '@angular/http';
import { env } from '../../../environments/environment';
import * as matchInfo from 'matchInfo';
import 'rxjs/add/operator/map';

@Component({
  selector: 'app-match',
  templateUrl: './match.component.html',
  styleUrls: ['./match.component.css']
})
export class MatchComponent implements OnInit {

  playersId: string[];

  matches: matchInfo.MatchInfo['matches'];
  opponentNames: matchInfo.Participant['userName'][];
  opponentAvatars: matchInfo.Participant['avatarURL'][];
  lastPlayedWords: string[];
  usedWords: string[][];
  letterGrids: matchInfo.Match['letters'][];
  tileGrids: matchInfo.Tile[][];

  selectedTile: boolean[][];
  foundWords: string[][];

  choosingWord: string[];

  constructor(private http: Http) {
    this.playersId = env.player.map(p => p.id);
  }

  ngOnInit() {
    this.fetchGames();
    window.onscroll = () => {
      if (window.scrollY == 0) {
        this.fetchGames();
        setTimeout(() => {
          if (window.scrollY == 0) window.scroll(0,99999);
        }, 800);
      }
    };
  }

  fetchGames() {
    this.http.get('http://' + window.location.host + '/match')
      .map((resp) => resp.text().length > 1000 ? resp.json() : '')
      .subscribe(
        (data) => {
          if (data) {
            this.processGameData(data);
            // auto find words for matches that already started
            for (let i = 0; i < this.matches.length; i++) {
              this.selectedTile[i] = Array<boolean>(25);
              if (this.matches[i].matchStatus === 4) {
                // For fresh new match, build heatmap
                this.buildHeatmap(i);
                continue; //matchStatus==4: new game, escape
              }
              this.selectAllWhite(i);
              this.findWords(i).then(isFound => {
                if (!isFound) this.selectAllPink(i);
              })
            }
          }
        });
  }

  processGameData(data) {
    if (data.matches) {
      // fetch exsisting game list
      this.matches = data.matches;
    } else if (data.match) {
      // fetch newly created game
      this.matches = [data.match];
    } else {
      console.log('Cannot fetch match data from MITM server...');
      return;
    }

    // only show matches on my turn
    this.matches = this.matches.filter(m => this.playersId.includes(m.participants[m.currentPlayerIndex].userId));

    // sort new turns on top
    // this.matches.sort((b, a) => a.participants[a.currentPlayerIndex==1?0:1].turnDate.valueOf()
    //   - b.participants[b.currentPlayerIndex==1?0:1].turnDate.valueOf());
    console.log(this.matches);

    const opponents = this.matches.map(m => this.playersId.includes(m.participants[0].userId) ? m.participants[1] : m.participants[0]);
    this.opponentNames = opponents.map(p => p.userName);
    this.opponentAvatars = opponents.map(p => p.avatarURL ? p.avatarURL : 'https://thesocietypages.org/socimages/files/2009/05/nopic_192.gif');
    this.usedWords = this.matches.map(m => m.serverData.usedWords);
    this.lastPlayedWords = this.usedWords.map(ws => ws.length > 0 ? ws[ws.length - 1] : '');

    this.letterGrids = this.matches.map(m => m['letters']);
    console.log(this.letterGrids);
    this.tileGrids = this.matches.map(m => m['serverData']['tiles'])
      .map(t => t.slice(20, 25).concat(t.slice(15, 20).concat(t.slice(10, 15).concat(t.slice(5, 10).concat(t.slice(0, 5))))));
    this.tileGrids.forEach(tg => tg.forEach(t => t.t = t.t.toUpperCase()));

    // build tileGrids
    for (let i = 0; i < this.tileGrids.length; i++) {
      const tg = this.tileGrids[i];
      // reverse owner, if player moves first
      if (this.playersId.includes(this.matches[i].participants[0].userId)) {
        tg.forEach(t => t.o == 1 ? t.o = 0 : t.o == 0 ? t.o = 1 : t.o)
      }

      // set surronded tiles: if every surrounded tiles have same ownership, set true; otherwise, false
      for (let k = 0; k < 25; k++) {
        tg[k].s = [
          ([4, 9, 14, 19, 24].includes(k)) ? undefined : tg[k + 1], // right
          ([0, 5, 10, 15, 20].includes(k)) ? undefined : tg[k - 1], // left
          tg[k + 5], // down
          tg[k - 5]  // up
        ]
          .filter(t => t)
          .every(t => t.o === tg[k].o); // tg[k].o: ownership of adjacent tiles of t; t.o: ownership of t

        tg[k].color = tg[k].o === 127 ? 'white'
          : tg[k].o === 0 && tg[k].s ? 'red'
            : tg[k].o === 0 && !tg[k].s ? 'pink'
              : tg[k].o === 1 && tg[k].s ? 'blue'
                : tg[k].o === 1 && !tg[k].s ? 'azure' : 'error';
      }
    }

    this.selectedTile = Array<boolean[]>(this.matches.length);
    this.foundWords = Array<string[]>(this.matches.length);
    this.choosingWord = Array<string>(this.matches.length);
  }

  //TODO: Remove the recursive call. The second param decides if the recursive call going on. 
  findWords(i: number): Promise<boolean> {
    const letters = this.tileGrids[i].map(t => t.t).join('').toUpperCase();
    const tg = this.tileGrids[i]; 
    // Handicap to weak opponent
    let numSelected = this.selectedTile[i].filter(t => t).length;
    const numRed = tg.filter(t => t.color === 'red').length;
    const numPink = tg.filter(t => t.color === 'pink').length;
    const numBlue = tg.filter(t => t.color === 'blue').length
    const numAzure = tg.filter(t => t.color === 'azure').length;
    while (numSelected > 1 && (numSelected+numAzure+numBlue/3) / (numPink+numRed) > 2 && (numBlue > numRed || numRed===0)) {
      if(!this.unselectLetter(i, tg)) break;
      numSelected = this.selectedTile[i].filter(t => t).length;
    }

    // Build selected letters
    let selected = [];
    for (let k = 0; k < 25; k++) {
      if (this.selectedTile[i][k]) {
        selected.push(letters[k])
      }
    }
    console.log(letters, selected.join(''));
    return new Promise(resolve => {
     this.http.get('http://' + window.location.host + '/words?selected=' + selected.join('') + '&letters=' + letters)
      .map(resp => resp.json())
      .subscribe(data => {
        this.foundWords[i] = data;
        const usedWords = this.usedWords[i];
        // filter out usedWords nor words as begining of usedWords
        this.foundWords[i] = this.foundWords[i].filter(w => !usedWords.some(uw => uw.indexOf(w.replace('*', '')) === 0));
        
        const existSelectedNonBlank = this.selectedTile[i].map((b, k) => tg[k].o!==127 && b).some(b => b);
        if (this.foundWords[i].length == 0) {
          // Test whether all selected tiles are blank
          if (!existSelectedNonBlank) {
            resolve(false);
            return;
          }
          // Recursively unselect letters
          this.unselectLetter(i, tg)
          this.findWords(i);
        } else {
          // Try to move cursor to the next virgin word (never played)
          const virginIndex = this.foundWords[i].findIndex(w => w.charAt(w.length-1) !== '*')
          if (virginIndex < 0) {
            this.unselectLetter(i, tg);
            if (existSelectedNonBlank && numBlue > numRed) {
              this.findWords(i);
            }
          }
          this.choosingWord[i] = this.foundWords[i][Math.max(0, virginIndex)];
        }
        
        resolve(true);
      });
    });
  }

  // Set colorCode for each tile based on all available words hits rarity
  buildHeatmap(i: number) {
    const letters = this.tileGrids[i].map(t => t.t).join('').toUpperCase();
    this.http.get('http://' + window.location.host + '/letterFrequency?letters=' + letters)
    .map((resp) => resp.json())
    .subscribe((data :number[]) => {
      const alphabet = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
      const freq = {};
      // Data::int[27] wC,A,...,Z
      const wordCount = data.shift();
      data.forEach((d,i) => { if (d > 0) freq[alphabet.charAt(i)] = d; });
      const maxFreq = Math.max.apply(null, data);
      const minFreq = Math.min.apply(null, data.filter(d => d>0));
      const freqList = Object.keys(freq).map(k => [k, freq[k]]).sort((a,b) => a[1]-b[1]);
      console.log('letter freq(A-Z)', freqList);
      console.log('word count', wordCount);
      console.log('max frequency', maxFreq);
      console.log('min frequency', minFreq);
      const tg = this.tileGrids[i];
      for (let k = 0; k < 25; k++) {
        tg[k].colorCode = `hsl(103, 90%, ${(Math.log1p(freq[tg[k].t]) / Math.log1p(maxFreq)) ** 1 * 100}%`;
      }
      speechSynthesis.speak(new SpeechSynthesisUtterance(Math.round(wordCount/1000)+'K;'));
    })
  }

  unselectLetter(i: number, tg: matchInfo.Tile[]): boolean {
    if (this.selectedTile[i].filter(t => t).length === 0) return false;

    const positionOrder = [[12], [11, 13, 7, 17], [6, 8, 16, 18], [2, 10, 14, 22], [1, 3, 5, 9, 15, 19, 21, 23], [0, 4, 20, 24]];
    const frequencyOrder = 'ESIARNOTLCDUPMGHBYFVKWZXQJ';
    let kToUnselect = -1;
    for (let ps of positionOrder) {
      for (let q = 0; q < frequencyOrder.length; q++) {
        const l = frequencyOrder[q];
        for (let k of ps) {
          if(this.selectedTile[i][k] && tg[k].t === l && tg[k].o !== 127){
            kToUnselect = k;
            break;
          }
        }
        if (kToUnselect >=0) break;
      }
      if (kToUnselect >=0) break;
    }
    if (kToUnselect < 0) return false;
    this.selectedTile[i][kToUnselect] = false;
    return true;
  }

  clearSelected(i: number) {
    this.selectedTile[i].fill(false)
  }

  // select all untouched tiles (white)
  selectAllWhite(i: number) {
    const tg = this.tileGrids[i];
    for (let k = 0; k < 25; k++) {
      this.selectedTile[i][k] = (tg[k].o === 127);
    }
  }

  // select all unsurrounded opponent's tiles (pink)
  selectAllPink(i: number) {
    const tg = this.tileGrids[i];
    for (let k = 0; k < 25; k++) {
      this.selectedTile[i][k] = (tg[k].o == 0 && !tg[k].s);
    } 
    this.findWords(i);
  }

  autoClickTiles(i: number) {
    if (this.choosingWord[i] === undefined) return;

    const tg = this.tileGrids[i];
    const letterMap = {}; // key: letter; value: positions[]

    // construct the letterMap, each letter key will have ordered (pink,white...) tiles number as value
    // 1st, push selected tiles
    for (let k = 0; k < 25; k++) {
      if (this.selectedTile[i][k]) {
        if (!letterMap[tg[k].t]) {
          letterMap[tg[k].t] = []
        }
        letterMap[tg[k].t].push(k)
      }
    }
    // 2nd, push unselected tiles, pink first, then white
    const positionOrder = [0, 4, 20, 24, 1, 3, 5, 9, 15, 19, 21, 23, 2, 10, 14, 22, 6, 8, 16, 18, 11, 13, 7, 17, 12];
    ["pink", "white", "red", "azure", "blue"].forEach(color => {
      for (let k of positionOrder) {
        if (tg[k].color === color && !this.selectedTile[i][k]) {
          if (!letterMap[tg[k].t]) {
            letterMap[tg[k].t] = []
          }
          letterMap[tg[k].t].push(k)
        }
      }
    })

    // decide the order of tiles to be clicked
    const clickOrderList = []
    for (let idx = 0; idx < this.choosingWord[i].length; idx++) {
      const letter = this.choosingWord[i][idx].toUpperCase();
      if (/^[A-Z]$/.test(letter)) {
        clickOrderList.push(letterMap[letter].shift())
      }
    }

    this.http.post('http://' + window.location.host + '/word?click=' + this.choosingWord[i], clickOrderList)
      .subscribe()
    console.log('autoclicking tiles...', this.choosingWord[i]);
  }

  deleteWord(i: number) {
    this.http.delete('http://' + window.location.host + '/word?delete=' + this.choosingWord[i])
    .subscribe()
    console.log('deleting...', this.choosingWord[i]);
  }
  
  // addWords(lastPlayedWords) {
  //   setTimeout(() => {      
  //     this.http.put('http://' + window.location.host + '/words', lastPlayedWords)
  //     .subscribe()
  //   }, 3000);
  // }

  speakWordUS(word: string) {
    word = word.replace(/\W/gi, '');
    let s = new SpeechSynthesisUtterance(word);
    s.voice = speechSynthesis.getVoices().filter(v => v.lang.indexOf('en-US') >= 0)[19]
    speechSynthesis.speak(s);
  }
  speakWordUK(word: string) {
    word = word.replace(/\W/gi, '');
    let s = new SpeechSynthesisUtterance(word);
    s.voice = speechSynthesis.getVoices().filter(v => v.lang.indexOf('en-GB') >= 0)[1]
    speechSynthesis.speak(s);
  }
  speakWordFR(word: string) {
    word = word.replace(/\W/gi, '');
    let s = new SpeechSynthesisUtterance(word);
    s.voice = speechSynthesis.getVoices().filter(v => v.lang.indexOf('fr-FR') >= 0)[1]
    speechSynthesis.speak(s);
  }
  speakWordNL(word: string) {
    word = word.replace(/\W/gi, '');
    let s = new SpeechSynthesisUtterance(word);
    s.voice = speechSynthesis.getVoices().filter(v => v.lang.indexOf('nl-NL') >= 0)[1]
    speechSynthesis.speak(s);
  }
  speakWordDE(word: string) {
    word = word.replace(/\W/gi, '');
    let s = new SpeechSynthesisUtterance(word);
    s.voice = speechSynthesis.getVoices().filter(v => v.lang.indexOf('de-DE') >= 0)[1]
    speechSynthesis.speak(s);
  }
}
