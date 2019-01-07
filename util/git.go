package util

import (
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// store h in set, s, handling nil s if necessary. Return new set.
func store(s map[plumbing.Hash]bool, h plumbing.Hash) map[plumbing.Hash]bool {
	if s == nil {
		s = make(map[plumbing.Hash]bool)
	}
	s[h] = true
	return s
}

// mergeBase finds best common ancestors between two commits to use in a
// three-way merge. One common ancestor is better than another common ancestor
// if the latter is an ancestor of the former. A common ancestor that does not
// have any better common ancestor is a best common ancestor, i.e. a merge base.
// Note that there can be more than one merge base for a pair of commits.
func MergeBase(repo *git.Repository, a, b plumbing.Hash) ([]plumbing.Hash, error) {
	s := repo.Storer
	commitA, err := object.GetCommit(s, a)
	if err != nil {
		return nil, err
	}

	commitB, err := object.GetCommit(s, b)
	if err != nil {
		return nil, err
	}

	// Mapping of direct descendants of each commit we visit
	desc := make(map[plumbing.Hash]map[plumbing.Hash]bool)

	// Set of commits reachable from a
	reachableFromA := make(map[plumbing.Hash]bool)

	// Walk commits reachable from A
	err = object.NewCommitPreorderIter(commitA, nil, nil).ForEach(func(c *object.Commit) error {
		reachableFromA[c.Hash] = true
		for _, h := range c.ParentHashes {
			desc[h] = store(desc[h], c.Hash)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Set of common commits between a and b
	common := make(map[plumbing.Hash]bool)

	// Walk commits reachable from B
	err = object.NewCommitPreorderIter(commitB, nil, nil).ForEach(func(c *object.Commit) error {
		if reachableFromA[c.Hash] {
			common[c.Hash] = true
		}
		for _, h := range c.ParentHashes {
			desc[h] = store(desc[h], c.Hash)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	best := make(map[plumbing.Hash]bool)

	// Trim down the set of common commits to only those that are best
	for h := range common {
		best[h] = true
		for child := range desc[h] {
			if common[child] {
				// there is a descendant to h that is common to both a and b. h is not in best.
				delete(best, h)
				break
			}
		}
	}

	var result []plumbing.Hash
	for h := range best {
		result = append(result, h)
	}
	return result, nil
}
