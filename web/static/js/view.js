var DateTime = luxon.DateTime;

function fill(data) {
  let repos = {};
  for (items of Object.values(data.data)) {
    for (i of Object.values(items)) {
      repos[i.id.owner+'/'+i.id.name] = i;
    }
  }

  // Use default icons
  document.querySelector(".octicon-star").id = 'icon-star';

  items = document.getElementsByTagName('li');
  for (item of items) {
    link = item.firstElementChild;
    if (link === undefined || link === null) {
      continue
    }
    res = /https?:\/\/github.com\/([^\/]*)\/([^\/]*)\/?/.exec(link.href);
    if ((res === null) || (res.length !== 3)) {
      continue
    }
    let repoID = res[1] + '/' + res[2];
    if (!(repoID in repos)) {
      continue
    }
    let box = document.createElement("div");
    item.prepend(box);

    // Add nodes.
    let updated_at = DateTime.fromISO(repos[repoID].last_commit);
    box.innerHTML = [
      '<div class="awg star">',
      '<svg width=16px height=16px viewBox="0 0 16 16">',
      '<use href="#icon-star"></use></svg>',
      '<span>',
      repos[repoID].star,
      '</span>',
      '</div>',
      '<div class="awg updated-at">',
      repos[repoID].last_commit,
      '</div>',
      '<div class="awg rel-updated-at">',
      updated_at.toRelativeCalendar(),
      '</div>',
    ].join('\n');
  }

  let style = document.createElement('style');
  style.innerHTML = `
  .awg {
    display: inline-block;
    font-size: .8rem;
    border-radius: .5rem;
    margin-right: .5rem;
    padding: .1rem .4rem 0 .4rem;
    color: #fff;
    background: #78c878;
  }
  .awg.star {
    min-width: 4rem;
  }
  .awg.star svg {
    vertical-align: -.25rem;
  }
  `
  let ref = document.querySelector('script');
  ref.parentNode.insertBefore(style, ref);
}

function main() {
  fetch(window.data_url)
    .then(resp => resp.json())
    .then(data => {
      fill(data);
    });
}

main();
